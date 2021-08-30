package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gravitytypes "github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	tendermintjson "github.com/tendermint/tendermint/libs/json"
)

func withPristineE2EEnvironment(t *testing.T, cb func(
	string,
	*dockertest.Pool,
	*dockertest.Network,
)) {
	t.Helper()

	chain := Chain{
		DataDir: "testdata",
		ID:      "testchain",
	}

	err := chain.CreateAndInitializeValidators(4)
	require.NoError(t, err, "error initializing validators")

	err = chain.CreateAndInitializeOrchestrators(uint8(len(chain.Validators)))
	require.NoError(t, err, "error initializing orchestrators")

	// add validator accounts to genesis file
	configDir := chain.Validators[0].ConfigDir()
	for _, n := range chain.Validators {
		err = addGenesisAccount(configDir, "", n.KeyInfo.GetAddress(), "100000000000stake,100000000000footoken")
		require.NoError(t, err, "error creating validator accounts")
	}

	// add orchestrator accounts to genesis file
	for _, n := range chain.Orchestrators {
		err = addGenesisAccount(configDir, "", n.KeyInfo.GetAddress(), "100000000000stake,100000000000footoken")
		require.NoError(t, err, "error creating orchestrator accounts")
	}

	// file_copy around the genesis file with the accounts
	for _, v := range chain.Validators[1:] {
		_, err = fileCopy(filepath.Join(configDir, "config", "genesis.json"), filepath.Join(v.ConfigDir(), "config", "genesis.json"))
		require.NoError(t, err, "error copying over genesis files")
	}

	// generate ethereum keys for validators,
	// add them to the ethereum genesis
	ethGenesis := EthereumGenesis{
		Difficulty: "0x400",
		GasLimit:   "0xB71B00",
		Config:     EthereumConfig{ChainID: 15},
		Alloc:      make(map[string]Allocation, len(chain.Validators)+1),
	}
	ethGenesis.Alloc["0xBf660843528035a5A4921534E156a27e64B231fE"] = Allocation{Balance: "0x1337000000000000000000"}
	for _, v := range chain.Validators {
		err = v.generateEthereumKey()
		require.NoError(t, err, "error copying over genesis files")

		ethGenesis.Alloc[v.EthereumKey.Address] = Allocation{Balance: "0x1337000000000000000000"}
	}

	// write out the genesis file
	ethGenesisMarshal, err := json.MarshalIndent(ethGenesis, "", "  ")
	require.NoError(t, err, "error marshalling ethereum genesis file")

	writeFile(t, filepath.Join(chain.ConfigDir(), "ETHGenesis.json"), ethGenesisMarshal)

	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(configDir)
	config.Moniker = chain.Validators[0].Moniker

	genFilePath := config.GenesisFile()
	appState, genDoc, err := types.GenesisStateFromGenFile(genFilePath)
	require.NoError(t, err, "error reading genesis file")

	var bank Bank
	err = json.Unmarshal(appState["bank"], &bank)
	require.NoError(t, err, "error unmarshalling bank genesis state")
	bank.DenomMetadata = append(bank.DenomMetadata, DenomMetadata{
		Description: "footoken",
		Display:     "mfootoken",
		Base:        "footoken",
		DenomUnits: []DenomUnit{
			{
				Denom:    "footoken",
				Exponent: 0,
			},
			{
				Denom:    "mfootoken",
				Exponent: 6,
			},
		},
	})
	bank.DenomMetadata = append(bank.DenomMetadata, DenomMetadata{
		Description: "stake",
		Display:     "mstake",
		Base:        "stake",
		DenomUnits: []DenomUnit{
			{
				Denom:    "stake",
				Exponent: 0,
			},
			{
				Denom:    "mstake",
				Exponent: 3,
			},
		},
	})

	bz, err := json.Marshal(bank)
	require.NoError(t, err, "error marshalling bank state")
	appState["bank"] = bz

	var genUtil GenUtil
	err = json.Unmarshal(appState["genutil"], &genUtil)
	require.NoError(t, err, "error unmarshalling genesis state")

	// generate gentxs
	amount, _ := sdktypes.NewIntFromString("100000000000")
	coin := sdktypes.Coin{Denom: "stake", Amount: amount}
	genTxs := make([]json.RawMessage, len(chain.Validators))

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdktypes.Msg)(nil), &stakingtypes.MsgCreateValidator{}, &gravitytypes.MsgDelegateKeys{})
	interfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{}, &ed25519.PubKey{})
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	for i, v := range chain.Validators {
		cvm, err := v.buildCreateValidatorMsg(coin)
		require.NoError(t, err, "error building create validator msg")

		dm := v.buildDelegateKeysMsg()
		require.NoError(t, err, "error building delegate keys msg")

		signedTx, err := v.signMsg(cvm, dm)
		require.NoError(t, err, "error signing create validator msg")

		txRaw, err := marshaler.MarshalJSON(signedTx)
		require.NoError(t, err, "error marshalling tx")

		genTxs[i] = txRaw
	}
	genUtil.GenTxs = genTxs

	bz, err = json.Marshal(genUtil)
	require.NoError(t, err, "error marshalling gen_util state")
	appState["genutil"] = bz

	bz, err = json.Marshal(appState)
	require.NoError(t, err, "error marshalling app state")

	genDoc.AppState = bz
	out, err := tendermintjson.MarshalIndent(genDoc, "", "  ")
	require.NoError(t, err, "error marshalling genesis doc")

	for _, validator := range chain.Validators {
		writeFile(t, filepath.Join(validator.ConfigDir(), "config", "genesis.json"), out)
	}

	// update config.toml files
	for i, v := range chain.Validators {
		var configToml ValidatorConfig
		path := filepath.Join(v.ConfigDir(), "config", "config.toml")
		_, err = toml.DecodeFile(path, &configToml)
		require.NoError(t, err, "error decoding config toml")

		configToml.P2P.Laddr = "tcp://0.0.0.0:26656"
		configToml.P2P.AddrBookStrict = false
		configToml.P2P.ExternalAddress = fmt.Sprintf("%s:%d", v.instanceName(), 26656)
		configToml.RPC.Laddr = "tcp://0.0.0.0:26657"
		configToml.StateSync.Enable = false

		if i > 0 {
			configToml.LogLevel = "info"
		}

		var peers []string

		for j := 0; j < len(chain.Validators); j++ {
			if i == j {
				continue
			}
			peer := chain.Validators[j]
			peerID := fmt.Sprintf("%s@%s%d:26656", peer.NodeKey.ID(), peer.Moniker, j)
			peers = append(peers, peerID)
		}

		configToml.P2P.PersistentPeers = strings.Join(peers, ",")

		var b bytes.Buffer
		encoder := toml.NewEncoder(&b)
		err = encoder.Encode(configToml)
		require.NoError(t, err, "error encoding config toml")

		// todo(levi) use writeFile?
		err = os.WriteFile(path, b.Bytes(), fs.ModePerm)
		require.NoError(t, err, "error writing config toml")
	}

	// bring up docker network
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "error creating docker pool")
	network, err := pool.CreateNetwork("testnet")
	require.NoError(t, err, "error creating testnet network")

	// bring up ethereum
	t.Log("building and running ethereum")
	ethereum, err := pool.BuildAndRunWithBuildOptions(
		&dockertest.BuildOptions{
			Dockerfile: "ethereum/Dockerfile",
			ContextDir: "./",
		},
		&dockertest.RunOptions{
			Name:      "ethereum",
			NetworkID: network.Network.ID,
			PortBindings: map[docker.Port][]docker.PortBinding{
				"8545/tcp": {{HostIP: "", HostPort: "8545"}},
			},
			Env: []string{},
		},
		noRestart,
	)
	require.NoError(t, err, "error bringing up ethereum")
	t.Logf("deployed ethereum at %s", ethereum.Container.ID)

	wd, err := os.Getwd()
	require.NoError(t, err, "couldn't get working directory")

	for _, validator := range chain.Validators {
		runOpts := &dockertest.RunOptions{
			Name:       validator.instanceName(),
			NetworkID:  network.Network.ID,
			Mounts:     []string{fmt.Sprintf("%s/testdata/%s/%s/:/root/home", wd, chain.ID, validator.instanceName())},
			Repository: "gravity",
			Tag:        "prebuilt",
		}

		// expose the first validator for debugging and communication
		if validator.Index == 0 {
			runOpts.PortBindings = map[docker.Port][]docker.PortBinding{
				"1317/tcp":  {{HostIP: "", HostPort: "1317"}},
				"6060/tcp":  {{HostIP: "", HostPort: "6060"}},
				"6061/tcp":  {{HostIP: "", HostPort: "6061"}},
				"6062/tcp":  {{HostIP: "", HostPort: "6062"}},
				"6063/tcp":  {{HostIP: "", HostPort: "6063"}},
				"6064/tcp":  {{HostIP: "", HostPort: "6064"}},
				"6065/tcp":  {{HostIP: "", HostPort: "6065"}},
				"9090/tcp":  {{HostIP: "", HostPort: "9090"}},
				"26656/tcp": {{HostIP: "", HostPort: "26656"}},
				"26657/tcp": {{HostIP: "", HostPort: "26657"}},
			}
		}

		resource, err := pool.RunWithOptions(runOpts, noRestart)
		require.NoError(t, err, "error bringing up %s", validator.instanceName())
		t.Logf("deployed %s at %s", validator.instanceName(), resource.Container.ID)
	}

	// bring up the contract deployer and deploy contract
	t.Log("deploying contract_deployer")
	contractDeployer, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       "contract_deployer",
			Repository: "solidity",
			Tag:        "prebuilt",
			NetworkID:  network.Network.ID,
			PortBindings: map[docker.Port][]docker.PortBinding{
				"8545/tcp": {{HostIP: "", HostPort: "8545"}},
			},
			Env: []string{},
		},
		noRestart,
	)
	require.NoError(t, err, "error bringing up contract_deployer")
	t.Logf("deployed contract_deployer at %s", contractDeployer.Container.ID)

	container := contractDeployer.Container
	for container.State.Running {
		time.Sleep(10 * time.Second)
		container, err = pool.Client.InspectContainer(contractDeployer.Container.ID)
		require.NoError(t, err, "error inspecting contract deployer")
	}

	contractDeployerLogOutput := bytes.Buffer{}
	err = pool.Client.Logs(docker.LogsOptions{
		Container:    contractDeployer.Container.ID,
		OutputStream: &contractDeployerLogOutput,
		Stdout:       true,
	})
	require.NoError(t, err, "error getting contract deployer logs")

	var gravityContract string
	for _, s := range strings.Split(contractDeployerLogOutput.String(), "\n") {
		if strings.HasPrefix(s, "Gravity deployed at Address") {
			strSpl := strings.Split(s, "-")
			gravityContract = strings.ReplaceAll(strSpl[1], " ", "")
			break
		}
	}
	err = pool.RemoveContainerByName(container.Name)
	require.NoError(t, err, "error removing contract deployer container")
	require.NotEmptyf(t, gravityContract, "empty gravity contract")

	// deploy orchestrators
	for _, orchestrator := range chain.Orchestrators {
		validator := chain.Validators[orchestrator.Index]
		env := []string{
			fmt.Sprintf("VALIDATOR=%s", validator.instanceName()),
			fmt.Sprintf("COSMOS_GRPC=http://%s:9090/", validator.instanceName()),
			fmt.Sprintf("COSMOS_RPC=http://%s:1317", validator.instanceName()),
			fmt.Sprintf("VALIDATOR=%s", validator.instanceName()),
			fmt.Sprintf("COSMOS_PHRASE=%s", orchestrator.Mnemonic),
			fmt.Sprintf("ETH_PRIVATE_KEY=%s", validator.EthereumKey.PrivateKey),
			fmt.Sprintf("CONTRACT_ADDR=%s", gravityContract),
			"DENOM=stake",
			"ETH_RPC=http://ethereum:8545",
			"RUST_BACKTRACE=full",
		}
		runOpts := &dockertest.RunOptions{
			Name:       orchestrator.instanceName(),
			NetworkID:  network.Network.ID,
			Repository: "orchestrator",
			Tag:        "prebuilt",
			Env:        env,
		}

		resource, err := pool.RunWithOptions(runOpts, noRestart)
		require.NoError(t, err, "error bringing up %s", orchestrator.instanceName())
		t.Logf("deployed %s at %s", orchestrator.instanceName(), resource.Container.ID)
	}

	// write test runner files to config directory
	var ethKeys []string
	var validatorPhrases []string
	for _, validator := range chain.Validators {
		ethKeys = append(ethKeys, validator.EthereumKey.PrivateKey)
		validatorPhrases = append(validatorPhrases, validator.Mnemonic)
	}
	var orchestratorPhrases []string
	for _, orchestrator := range chain.Orchestrators {
		orchestratorPhrases = append(orchestratorPhrases, orchestrator.Mnemonic)
	}

	writeFile(t, filepath.Join(chain.DataDir, "validator-eth-keys"), []byte(strings.Join(ethKeys, "\n")))
	writeFile(t, filepath.Join(chain.DataDir, "validator-phrases"), []byte(strings.Join(validatorPhrases, "\n")))
	writeFile(t, filepath.Join(chain.DataDir, "orchestrator-phrases"), []byte(strings.Join(orchestratorPhrases, "\n")))
	writeFile(t, filepath.Join(chain.DataDir, "contracts"), contractDeployerLogOutput.Bytes())

	cb(wd, pool, network)
}

func noRestart(config *docker.HostConfig) {
	// in this case we don't want the nodes to restart on failure
	config.RestartPolicy = docker.RestartPolicy{
		Name: "no",
	}
}
