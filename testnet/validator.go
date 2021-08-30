package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/codec/unknownproto"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/go-bip39"
	"github.com/cosmos/gravity-bridge/module/app"
	gravitytypes "github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	cfg "github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Validator struct {
	Chain   *Chain
	Index   uint8
	Moniker string

	// Key management
	Mnemonic         string
	KeyInfo          keyring.Info
	PrivateKey       cryptotypes.PrivKey
	ConsensusKey     privval.FilePVKey
	ConsensusPrivKey cryptotypes.PrivKey
	NodeKey          p2p.NodeKey

	EthereumKey EthereumKey
}

type EthereumKey struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Address    string `json:"address"`
}

// createMnemonic creates a new mnemonic
func createMnemonic() (string, error) {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

func (v *Validator) ConfigDir() string {
	return fmt.Sprintf("%s/%s", v.Chain.ConfigDir(), v.instanceName())
}

// MkDir creates the directory for the testnode
func (v *Validator) MkDir() {
	p := path.Join(v.ConfigDir(), "config")
	if err := os.MkdirAll(p, 0755); err != nil {
		panic(err)
	}
}

func getGenDoc(path string) (doc *tmtypes.GenesisDoc, err error) {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)

	genFile := config.GenesisFile()
	doc = &tmtypes.GenesisDoc{}
	if _, err = os.Stat(genFile); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = nil
	} else {
		doc, err = tmtypes.GenesisDocFromFile(genFile)
		if err != nil {
			err = errors.Wrap(err, "Failed to read genesis doc from file")
			return
		}
	}

	return
}

func (v *Validator) init() error {
	encodingConfig := app.MakeEncodingConfig()
	cdc := encodingConfig.Marshaler

	v.MkDir()
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(v.ConfigDir())
	config.Moniker = v.Moniker

	genDoc, err := getGenDoc(v.ConfigDir())
	if err != nil {
		return err
	}

	appState, err := json.MarshalIndent(app.ModuleBasics.DefaultGenesis(cdc), "", " ")
	if err != nil {
		return errors.Wrap(err, "Failed to marshall default genesis state")
	}

	genDoc.ChainID = v.Chain.ID
	genDoc.Validators = nil
	genDoc.AppState = appState
	if err = genutil.ExportGenesisFile(genDoc, config.GenesisFile()); err != nil {
		return errors.Wrap(err, "Failed to export genesis file")
	}

	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
	return nil
}

func (v *Validator) createNodeKey() error {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(v.ConfigDir())
	config.Moniker = v.Moniker

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return err
	}
	v.NodeKey = *nodeKey
	return nil
}

func (v *Validator) createConsensusKey() (err error) {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(v.ConfigDir())
	config.Moniker = v.Moniker

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := tmos.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
		return err
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := tmos.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
		return err
	}

	filePV := privval.LoadOrGenFilePV(pvKeyFile, pvStateFile)
	v.ConsensusKey = filePV.Key
	return nil
}

// createMemoryKey creates a key but doesn't store it to any files
func createMemoryKey() (mnemonic string, info *keyring.Info, err error) {
	// Get bip39 mnemonic
	mnemonic, err = createMnemonic()
	if err != nil {
		return
	}

	account, err := createMemoryKeyFromMnemonic(mnemonic)
	return mnemonic, account, err
}

// createMemoryKey creates a key but doesn't store it to any files
func createMemoryKeyFromMnemonic(mnemonic string) (info *keyring.Info, err error) {
	kb, err := keyring.New("testnet", keyring.BackendMemory, "", nil)
	if err != nil {
		return
	}

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return
	}

	account, err := kb.NewAccount("", mnemonic, "", sdktypes.FullFundraiserPath, algo)
	info = &account
	return
}

func (v *Validator) createKeyFromMnemonic(name string, mnemonic string) (err error) {
	kb, err := keyring.New("testnet", keyring.BackendTest, v.ConfigDir(), nil)
	if err != nil {
		return err
	}

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return err
	}

	v.Mnemonic = mnemonic

	info, err := kb.NewAccount(name, mnemonic, "", sdktypes.FullFundraiserPath, algo)
	if err != nil {
		return err
	}
	v.KeyInfo = info

	privKeyArmor, err := kb.ExportPrivKeyArmor(name, "testpassphrase")
	if err != nil {
		return err
	}
	privKey, _, err := sdkcrypto.UnarmorDecryptPrivKey(privKeyArmor, "testpassphrase")
	if err != nil {
		return err
	}
	v.PrivateKey = privKey

	return nil
}

// createKey creates a new account and writes it to a validator's config directory
func (v *Validator) createKey(name string) (err error) {
	// Get bip39 mnemonic
	mnemonic, err := createMnemonic()
	if err != nil {
		return err
	}
	return v.createKeyFromMnemonic(name, mnemonic)
}

func addGenesisAccount(path string, moniker string, accAddr sdktypes.AccAddress, coinsStr string) (err error) {
	encodingConfig := app.MakeEncodingConfig()
	cdc := encodingConfig.Marshaler

	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)
	config.Moniker = moniker

	coins, err := sdktypes.ParseCoinsNormalized(coinsStr)
	if err != nil {
		return fmt.Errorf("failed to parse coins: %w", err)
	}

	balances := banktypes.Balance{Address: accAddr.String(), Coins: coins.Sort()}
	genAccount := authtypes.NewBaseAccount(accAddr, nil, 0, 0)

	genFile := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return fmt.Errorf("failed to get accounts from any: %w", err)
	}

	if accs.Contains(accAddr) {
		return fmt.Errorf("cannot add account at existing address %s", accAddr)
	}

	// Add the new account to the set of genesis accounts and sanitize the
	// accounts afterwards.
	accs = append(accs, genAccount)
	accs = authtypes.SanitizeGenesisAccounts(accs)

	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return fmt.Errorf("failed to convert accounts into any's: %w", err)
	}
	authGenState.Accounts = genAccs

	authGenStateBz, err := cdc.MarshalJSON(&authGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal auth genesis state: %w", err)
	}

	appState[authtypes.ModuleName] = authGenStateBz

	bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
	bankGenState.Balances = append(bankGenState.Balances, balances)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}

	appState[banktypes.ModuleName] = bankGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}

	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}

func (v *Validator) generateEthereumKey() (err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	v.EthereumKey = EthereumKey{
		PrivateKey: hexutil.Encode(privateKeyBytes),
		PublicKey:  hexutil.Encode(publicKeyBytes),
		Address:    crypto.PubkeyToAddress(*publicKeyECDSA).Hex(),
	}
	return
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func (v *Validator) buildCreateValidatorMsg(amount sdktypes.Coin) (sdktypes.Msg, error) {
	description := types.NewDescription(
		v.Moniker,
		"",
		"",
		"",
		"",
	)

	commissionRates := types.CommissionRates{
		Rate:          sdktypes.MustNewDecFromStr("0.1"),
		MaxRate:       sdktypes.MustNewDecFromStr("0.2"),
		MaxChangeRate: sdktypes.MustNewDecFromStr("0.01"),
	}

	// get the initial validator min self delegation
	minSelfDelegation, _ := sdktypes.NewIntFromString("1")

	valPubKey, err := cryptocodec.FromTmPubKeyInterface(v.ConsensusKey.PubKey)
	if err != nil {
		return nil, err
	}

	msg, err := types.NewMsgCreateValidator(
		sdktypes.ValAddress(v.KeyInfo.GetAddress()), valPubKey, amount, description, commissionRates, minSelfDelegation,
	)
	return msg, err
}

func (v *Validator) buildDelegateKeysMsg() sdktypes.Msg {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdktypes.Msg)(nil), &types.MsgCreateValidator{}, &gravitytypes.MsgDelegateKeys{})
	interfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{}, &ed25519.PubKey{})
	marshaller := codec.NewProtoCodec(interfaceRegistry)

	privKeyBz, err := hexutil.Decode(v.EthereumKey.PrivateKey)
	if err != nil {
		panic(fmt.Sprintf("failed to HEX decode private key: %s", err))
	}

	privKey, err := crypto.ToECDSA(privKeyBz)
	if err != nil {
		panic(fmt.Sprintf("failed to convert private key: %s", err))
	}

	signMsg := gravitytypes.DelegateKeysSignMsg{
		ValidatorAddress: sdktypes.ValAddress(v.KeyInfo.GetAddress()).String(),
		Nonce:            0,
	}

	signMsgBz := marshaller.MustMarshal(&signMsg)
	hash := crypto.Keccak256Hash(signMsgBz).Bytes()
	ethSig, err := gravitytypes.NewEthereumSignature(hash, privKey)
	if err != nil {
		panic(fmt.Sprintf("failed to create Ethereum signature: %s", err))
	}

	return gravitytypes.NewMsgDelegateKeys(
		sdktypes.ValAddress(v.KeyInfo.GetAddress()),
		v.Chain.Orchestrators[v.Index].KeyInfo.GetAddress(),
		v.EthereumKey.Address,
		ethSig,
	)
}

func (v *Validator) instanceName() string {
	return fmt.Sprintf("%s%d", v.Moniker, v.Index)
}

func decodeTx(txBytes []byte) (*tx.Tx, error) {
	var raw tx.TxRaw
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdktypes.Msg)(nil), &types.MsgCreateValidator{}, &gravitytypes.MsgDelegateKeys{})
	interfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{}, &ed25519.PubKey{})
	marshaller := codec.NewProtoCodec(interfaceRegistry)

	// reject all unknown proto fields in the root TxRaw
	err := unknownproto.RejectUnknownFieldsStrict(txBytes, &raw, interfaceRegistry)
	if err != nil {
		return nil, errors.Wrap(errors.ErrTxDecode, err.Error())
	}

	if err := marshaller.Unmarshal(txBytes, &raw); err != nil {
		return nil, err
	}

	var body tx.TxBody

	if err := marshaller.Unmarshal(raw.BodyBytes, &body); err != nil {
		return nil, errors.Wrap(errors.ErrTxDecode, err.Error())
	}

	var authInfo tx.AuthInfo

	// reject all unknown proto fields in AuthInfo
	err = unknownproto.RejectUnknownFieldsStrict(raw.AuthInfoBytes, &authInfo, interfaceRegistry)
	if err != nil {
		return nil, errors.Wrap(errors.ErrTxDecode, err.Error())
	}

	if err := marshaller.Unmarshal(raw.AuthInfoBytes, &authInfo); err != nil {
		return nil, errors.Wrap(errors.ErrTxDecode, err.Error())
	}

	theTx := &tx.Tx{
		Body:       &body,
		AuthInfo:   &authInfo,
		Signatures: raw.Signatures,
	}

	return theTx, nil
}

func (v *Validator) signMsg(msgs ...sdktypes.Msg) (*tx.Tx, error) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdktypes.Msg)(nil), &types.MsgCreateValidator{}, &gravitytypes.MsgDelegateKeys{})
	interfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{}, &ed25519.PubKey{})
	marshaller := codec.NewProtoCodec(interfaceRegistry)

	signModes := []txsigning.SignMode{txsigning.SignMode_SIGN_MODE_DIRECT}
	txConfig := authtx.NewTxConfig(marshaller, signModes)
	txBuilder := txConfig.NewTxBuilder()

	if err := txBuilder.SetMsgs(msgs...); err != nil {
		return nil, err
	}

	txBuilder.SetMemo(fmt.Sprintf("%s@%s:26656", v.NodeKey.ID(), v.instanceName()))
	fees := sdktypes.Coins{sdktypes.Coin{}}
	txBuilder.SetFeeAmount(fees)
	txBuilder.SetGasLimit(200000)
	txBuilder.SetTimeoutHeight(0)

	signerData := authsigning.SignerData{
		ChainID:       v.Chain.ID,
		AccountNumber: 0,
		Sequence:      0,
	}

	// For SIGN_MODE_DIRECT, calling SetSignatures calls setSignerInfos on
	// TxBuilder under the hood, and SignerInfos is needed to generated the
	// sign bytes. This is the reason for setting SetSignatures here, with a
	// nil signature.
	//
	// Note: this line is not needed for SIGN_MODE_LEGACY_AMINO, but putting it
	// also doesn't affect its generated sign bytes, so for code's simplicity
	// sake, we put it here.
	sigData := txsigning.SingleSignatureData{
		SignMode:  txsigning.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}
	sig := txsigning.SignatureV2{
		PubKey:   v.KeyInfo.GetPubKey(),
		Data:     &sigData,
		Sequence: 0,
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	bytesToSign, err := txConfig.SignModeHandler().GetSignBytes(txsigning.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	// Sign those bytes
	sigBytes, err := v.PrivateKey.Sign(bytesToSign)
	if err != nil {
		return nil, err
	}

	// Construct the SignatureV2 struct
	sigData = txsigning.SingleSignatureData{
		SignMode:  txsigning.SignMode_SIGN_MODE_DIRECT,
		Signature: sigBytes,
	}
	sig = txsigning.SignatureV2{
		PubKey:   v.KeyInfo.GetPubKey(),
		Data:     &sigData,
		Sequence: 0,
	}
	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	signedTx := txBuilder.GetTx()

	txEncoder := authtx.DefaultTxEncoder()
	bz, err := txEncoder(signedTx)
	if err != nil {
		return nil, err
	}

	stdTx, err := decodeTx(bz)

	return stdTx, err
}
