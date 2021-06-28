package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
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
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	gravitytypes "github.com/cosmos/gravity-bridge/module/x/gravity/types"
	dt "github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	tmjson "github.com/tendermint/tendermint/libs/json"
)

func writeFile(path string, body []byte) error {
	if _, err := os.Create(path); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, body, 0644); err != nil {
		return err
	}

	return nil
}

var ValidatorMnemonics = []string{
	"indoor rail elder dwarf concert cry surprise ice exhaust sword million square dinosaur merry brother panther evolve usual excess when annual grid path saddle",
	"veteran early defense faith brand siege rent ceiling boring end before genuine fault shrug crater raccoon success paper grunt prefer transfer powder repeat boil",
	"wonder armed degree prosper swarm crouch uphold hair deputy hospital title group unusual baby asthma bean mouse hobby custom era success earth wonder clog",
	"person blast robust group uniform awesome poem prevent pear sugar audit curious video punch senior nest short song topple sick guitar unhappy state forum",
}

var OrchestratorMnemonics = []string{
	"dwarf armed ordinary industry drive normal radar laugh host garden banner gather bless source elegant gold myth rail leader describe sick actor clown fuel",
	"print screen fade apart galaxy major flame kick power initial false window spoil caught pink find anxiety prosper elder speak enough often mountain deny",
	"theory system waste zone best still auction drop bread flavor caught daughter tower elder normal movie swallow crack just best ball exist seven involve",
	"ship beef little food chase ostrich neutral lemon abandon nut balance nice outside crater creek radio witness tree bridge yard depend curve fever wife",
}

var ValidatorEthereumKeys = []EthereumKey{
	{
		PublicKey:  "0x0429e9b421a18927707d23d74ec17a18d70dd582c566a02f70f825ae1a9be933946ee6753d40fa467037bbb33afed2e2031ad2008461ee15caa634d505cd9898a7",
		PrivateKey: "0xd2ef9c9fd43ab2a652f2ef343c7c7748d639fe1bf060b542222b74ac2a427bcd",
		Address:    "0x1b30Ed4DEF9933ef5a597713Ca61d28734F67BCA",
	},
	{
		PublicKey:  "0x04e373425bb40df0ac023b30bf33b917141b8abe5897a39bfae5f4f32a71efa0f5e84d355553c8b75e7f2f38326b4748f23158d4d0b22727a4bc90d4da897712e2",
		PrivateKey: "0x8c9cd3be82e86bfa77661768d79f2a899651f1c952ee2541488d6ee074a77c10",
		Address:    "0xAdD1fB646a7e93351E4c68801FD3Df224BD75e6B",
	},
	{
		PublicKey:  "0x0417aeb52eec2e05c825e428fc43149973d35ac29bc1126a71ff64e19708ee71873497cdd78bf025bd1ffd8f6c853bb5806e84870c993fd246d851171af9d07668",
		PrivateKey: "0x19f772379632aa5542e4a1637008d1702a3bd44a1c329af981a60cafd30bf946",
		Address:    "0x8C1A34491101E8f00204BBaF1ab34604a08F0b9A",
	},
	{
		PublicKey:  "0x0407e5e679602583a9d2735c84e3475b1305f4c4b608cc3e5c17402938a34e75656f8fe2e55b91fde56ac893637594061edc90a0719be491381a3621289bfb6d93",
		PrivateKey: "0xda4a30ccf1b6804d1f139d1912583825703e50478f9fd10436a8f0924a569061",
		Address:    "0xBb283c611036702ba33edC4C479daA102CCdBA12",
	},
}

var MinerKey = EthereumKey{
	PublicKey:  "",
	PrivateKey: "0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7",
	Address:    "0xBf660843528035a5A4921534E156a27e64B231fE",
}

func TestBasicChainDeterministicKeys(t *testing.T) {
	err := os.RemoveAll("testdata/")
	require.NoError(t, err, "unable to reset testdata directory")

	chain := Chain{
		DataDir:    "testdata",
		ID:         "testchain",
		Validators: nil,
	}

	err = chain.CreateAndInitializeValidatorsWithMnemonics(4, ValidatorMnemonics)
	require.NoError(t, err, "error initializing validators")

	err = chain.CreateAndInitializeOrchestratorsWithMnemonics(uint8(len(chain.Validators)), OrchestratorMnemonics)
	require.NoError(t, err, "error initializing orchestrators")

	// add validator accounts to genesis file
	path := chain.Validators[0].ConfigDir()
	for _, n := range chain.Validators {
		err = addGenesisAccount(path, "", n.KeyInfo.GetAddress(), "100000000000stake,100000000000footoken")
		require.NoError(t, err, "error creating validator accounts")
	}

	// add orchestrator accounts to genesis file
	for _, n := range chain.Orchestrators {
		err = addGenesisAccount(path, "", n.KeyInfo.GetAddress(), "100000000000stake,100000000000footoken")
		require.NoError(t, err, "error creating orchestrator accounts")
	}

	// file_copy around the genesis file with the accounts
	for _, v := range chain.Validators[1:] {
		_, err = fileCopy(filepath.Join(path, "config", "genesis.json"), filepath.Join(v.ConfigDir(), "config", "genesis.json"))
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
	for i, v := range chain.Validators {
		v.EthereumKey = ValidatorEthereumKeys[i]
		ethGenesis.Alloc[v.EthereumKey.Address] = Allocation{Balance: "0x1337000000000000000000"}
	}

	// write out the genesis file
	ethGenesisMarshal, err := json.MarshalIndent(ethGenesis, "", "  ")
	require.NoError(t, err, "error marshalling ethereum genesis file")

	err = ioutil.WriteFile(filepath.Join(chain.ConfigDir(), "ETHGenesis.json"), ethGenesisMarshal, 0644)
	require.NoError(t, err, "error writing ethereum genesis file")

	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)
	config.Moniker = chain.Validators[0].Moniker

	genFilePath := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFilePath)
	require.NoError(t, err, "error reading genesis file")

	var genUtil GenUtil
	err = json.Unmarshal(appState["genutil"], &genUtil)
	require.NoError(t, err, "error unmarshalling genesis state")

	// generate gentxs
	amount, _ := sdktypes.NewIntFromString("100000000000")
	coin := sdktypes.Coin{Denom: "stake", Amount: amount}
	genTxs := make([]json.RawMessage, len(chain.Validators))

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdktypes.Msg)(nil), &types.MsgCreateValidator{}, &gravitytypes.MsgDelegateKeys{})
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

	bz, err := json.Marshal(genUtil)
	require.NoError(t, err, "error marshalling gen_util state")
	appState["genutil"] = bz

	bz, err = json.Marshal(appState)
	require.NoError(t, err, "error marshalling app state")

	genDoc.AppState = bz
	out, err := tmjson.MarshalIndent(genDoc, "", "  ")
	require.NoError(t, err, "error marshalling genesis doc")

	for _, validator := range chain.Validators {
		err = ioutil.WriteFile(filepath.Join(validator.ConfigDir(), "config", "genesis.json"), out, fs.ModePerm)
		require.NoError(t, err, "error writing out genesis file")
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

		err = os.WriteFile(path, b.Bytes(), fs.ModePerm)
		require.NoError(t, err, "error writing config toml")

		startupPath := filepath.Join(v.ConfigDir(), "startup.sh")
		err = os.WriteFile(startupPath, []byte(fmt.Sprintf("gravity --home home start --pruning=nothing > home.n%d.log", v.Index)), fs.ModePerm)
	}

	// bring up docker network
	pool, err := dt.NewPool("")
	require.NoError(t, err, "error creating docker pool")
	network, err := pool.CreateNetwork("testnet")
	require.NoError(t, err, "error creating testnet network")
	defer func() {
		network.Close()
	}()

	hostConfig := func(config *dc.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = dc.RestartPolicy{
			Name: "no",
		}
	}

	// bring up ethereum
	t.Log("building and running ethereum")
	ethereum, err := pool.BuildAndRunWithBuildOptions(&dt.BuildOptions{
		Dockerfile: "ethereum/Dockerfile",
		ContextDir: "./",
	},
		&dt.RunOptions{
			Name:      "ethereum",
			NetworkID: network.Network.ID,
			PortBindings: map[dc.Port][]dc.PortBinding{
				"8545/tcp": {{HostIP: "", HostPort: "8545"}},
			},
			Env: []string{},
		}, hostConfig)
	require.NoError(t, err, "error bringing up ethereum")
	t.Logf("deployed ethereum at %s", ethereum.Container.ID)
	defer func() {
		ethereum.Close()
	}()

	// build validators
	for _, validator := range chain.Validators {
		t.Logf("building %s", validator.instanceName())
		err = pool.Client.BuildImage(dc.BuildImageOptions{
			Name:         validator.instanceName(),
			Dockerfile:   "Dockerfile",
			ContextDir:   "./module",
			OutputStream: ioutil.Discard,
		})
		require.NoError(t, err, "error building %s", validator.instanceName())
		t.Logf("built %s", validator.instanceName())
	}

	wd, err := os.Getwd()
	require.NoError(t, err, "couldn't get working directory")

	for _, validator := range chain.Validators {
		runOpts := &dt.RunOptions{
			Name:       validator.instanceName(),
			NetworkID:  network.Network.ID,
			Mounts:     []string{fmt.Sprintf("%s/testdata/testchain/%s/:/root/home", wd, validator.instanceName())},
			Repository: validator.instanceName(),
		}

		// expose the first validator for debugging and communication
		if validator.Index == 0 {
			runOpts.PortBindings = map[dc.Port][]dc.PortBinding{
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

		resource, err := pool.RunWithOptions(runOpts, hostConfig)
		require.NoError(t, err, "error bringing up %s", validator.instanceName())
		t.Logf("deployed %s at %s", validator.instanceName(), resource.Container.ID)
		defer func() {
			resource.Close()
		}()
	}

	// bring up the contract deployer and deploy contract
	t.Log("building contract_deployer")
	contractDeployer, err := pool.BuildAndRunWithBuildOptions(
		&dt.BuildOptions{
			Dockerfile: "Dockerfile",
			ContextDir: "./solidity",
		},
		&dt.RunOptions{
			Name:      "contract_deployer",
			NetworkID: network.Network.ID,
			PortBindings: map[dc.Port][]dc.PortBinding{
				"8545/tcp": {{HostIP: "", HostPort: "8545"}},
			},
			Env: []string{},
		}, func(config *dc.HostConfig) {})
	require.NoError(t, err, "error bringing up contract deployer")
	t.Logf("deployed contract deployer at %s", contractDeployer.Container.ID)
	defer func() {
		contractDeployer.Close()
	}()

	container := contractDeployer.Container
	for container.State.Running {
		time.Sleep(10 * time.Second)
		container, err = pool.Client.InspectContainer(contractDeployer.Container.ID)
		require.NoError(t, err, "error inspecting contract deployer")
	}

	contractDeployerLogOutput := bytes.Buffer{}
	err = pool.Client.Logs(dc.LogsOptions{
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

	// build orchestrators
	for _, orchestrator := range chain.Orchestrators {
		t.Logf("building %s", orchestrator.instanceName())
		err = pool.Client.BuildImage(dc.BuildImageOptions{
			Name:         orchestrator.instanceName(),
			Dockerfile:   "Dockerfile",
			ContextDir:   "./orchestrator",
			OutputStream: ioutil.Discard,
		})
		require.NoError(t, err, "error building %s", orchestrator.instanceName())
		t.Logf("built %s", orchestrator.instanceName())
	}

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
		runOpts := &dt.RunOptions{
			Name:       orchestrator.instanceName(),
			NetworkID:  network.Network.ID,
			Repository: orchestrator.instanceName(),
			Env:        env,
		}

		resource, err := pool.RunWithOptions(runOpts, hostConfig)
		require.NoError(t, err, "error bringing up %s", orchestrator.instanceName())
		t.Logf("deployed %s at %s", orchestrator.instanceName(), resource.Container.ID)
		defer func() {
			resource.Close()
		}()
	}

	// distribute ethereum from miner to validators
	// todo: Implement happy path directly without the test_runner container
}
