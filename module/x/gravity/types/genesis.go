package types

import (
	"bytes"
	"fmt"
	"time"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
)

// DefaultParamspace defines the default auth module parameter subspace
const (
	// todo: implement oracle constants as params
	DefaultParamspace     = ModuleName
	EventVoteRecordPeriod = 24 * time.Hour // TODO: value????
)

var (
	// ParamsStoreKeyGravityID stores the gravity id
	ParamsStoreKeyGravityID = []byte("GravityID")

	// ParamsStoreKeyContractHash stores the contract hash
	ParamsStoreKeyContractHash = []byte("ContractHash")

	// ParamsStoreKeyBridgeContractAddress stores the contract address
	ParamsStoreKeyBridgeContractAddress = []byte("BridgeContractAddress")

	// ParamsStoreKeyBridgeContractChainID stores the bridge chain id
	ParamsStoreKeyBridgeContractChainID = []byte("BridgeChainID")

	// ParamsStoreKeySignedSignerSetTxsWindow stores the signed blocks window
	ParamsStoreKeySignedSignerSetTxsWindow = []byte("SignedSignerSetTxWindow")

	// ParamsStoreKeySignedBatchesWindow stores the signed blocks window
	ParamsStoreKeySignedBatchesWindow = []byte("SignedBatchesWindow")

	// ParamsStoreKeyEthereumSignaturesWindow stores the signed blocks window
	ParamsStoreKeyEthereumSignaturesWindow = []byte("EthereumSignaturesWindow")

	// ParamsStoreKeyTargetEthTxTimeout stores the target ethereum transaction timeout
	ParamsStoreKeyTargetEthTxTimeout = []byte("TargetEthTxTimeout")

	// ParamsStoreKeyAverageBlockTime stores the signed blocks window
	ParamsStoreKeyAverageBlockTime = []byte("AverageBlockTime")

	// ParamsStoreKeyAverageEthereumBlockTime stores the signed blocks window
	ParamsStoreKeyAverageEthereumBlockTime = []byte("AverageEthereumBlockTime")

	// ParamsStoreSlashFractionSignerSetTx stores the slash fraction valset
	ParamsStoreSlashFractionSignerSetTx = []byte("SlashFractionSignerSetTx")

	// ParamsStoreSlashFractionBatch stores the slash fraction Batch
	ParamsStoreSlashFractionBatch = []byte("SlashFractionBatch")

	// ParamsStoreSlashFractionEthereumSignature stores the slash fraction ethereum siganture
	ParamsStoreSlashFractionEthereumSignature = []byte("SlashFractionEthereumSignature")

	// ParamsStoreSlashFractionConflictingEthereumSignature stores the slash fraction ConflictingEthereumSignature
	ParamsStoreSlashFractionConflictingEthereumSignature = []byte("SlashFractionConflictingEthereumSignature")

	//  ParamStoreUnbondSlashingSignerSetTxsWindow stores unbond slashing valset window
	ParamStoreUnbondSlashingSignerSetTxsWindow = []byte("UnbondSlashingSignerSetTxsWindow")

	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{}
)

func (gs *GenesisState) UnpackInterfaces(unpacker cdctypes.AnyUnpacker) error {
	for _, otx := range gs.OutgoingTxs {
		var outgoing OutgoingTx
		if err := unpacker.UnpackAny(otx, &outgoing); err != nil {
			return err
		}
	}
	for _, sig := range gs.Confirmations {
		var signature EthereumTxConfirmation
		if err := unpacker.UnpackAny(sig, &signature); err != nil {
			return err
		}
	}
	for _, evr := range gs.EthereumEventVoteRecords {
		if err := evr.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}
	return nil
}

func EventVoteRecordPowerThreshold(totalPower sdk.Int) sdk.Int {
	return sdk.NewInt(66).Mul(totalPower).Quo(sdk.NewInt(100))
}

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (s GenesisState) ValidateBasic() error {
	if err := s.Params.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}
	return nil
}

// DefaultGenesisState returns empty genesis state
// TODO: set some better defaults here
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// DefaultParams returns a copy of the default params
func DefaultParams() *Params {
	return &Params{
		GravityId:                                 "defaultgravityid",
		BridgeEthereumAddress:                     "0x0000000000000000000000000000000000000000",
		SignedSignerSetTxsWindow:                  10000,
		SignedBatchesWindow:                       10000,
		EthereumSignaturesWindow:                  10000,
		TargetEthTxTimeout:                        43200000,
		AverageBlockTime:                          5000,
		AverageEthereumBlockTime:                  15000,
		SlashFractionSignerSetTx:                  sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionBatch:                        sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionEthereumSignature:            sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionConflictingEthereumSignature: sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		UnbondSlashingSignerSetTxsWindow:          10000,
	}
}

// ValidateBasic checks that the parameters have valid values.
func (p Params) ValidateBasic() error {
	if err := validateGravityID(p.GravityId); err != nil {
		return sdkerrors.Wrap(err, "gravity id")
	}
	if err := validateContractHash(p.ContractSourceHash); err != nil {
		return sdkerrors.Wrap(err, "contract hash")
	}
	if err := validateBridgeContractAddress(p.BridgeEthereumAddress); err != nil {
		return sdkerrors.Wrap(err, "bridge contract address")
	}
	if err := validateBridgeChainID(p.BridgeChainId); err != nil {
		return sdkerrors.Wrap(err, "bridge chain id")
	}
	if err := validateTargetEthTxTimeout(p.TargetEthTxTimeout); err != nil {
		return sdkerrors.Wrap(err, "Batch timeout")
	}
	if err := validateAverageBlockTime(p.AverageBlockTime); err != nil {
		return sdkerrors.Wrap(err, "Block time")
	}
	if err := validateAverageEthereumBlockTime(p.AverageEthereumBlockTime); err != nil {
		return sdkerrors.Wrap(err, "Ethereum block time")
	}
	if err := validateSignedSignerSetTxsWindow(p.SignedSignerSetTxsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSignedBatchesWindow(p.SignedBatchesWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateEthereumSignaturesWindow(p.EthereumSignaturesWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSlashFractionSignerSetTx(p.SlashFractionSignerSetTx); err != nil {
		return sdkerrors.Wrap(err, "slash fraction signersettx")
	}
	if err := validateSlashFractionBatch(p.SlashFractionBatch); err != nil {
		return sdkerrors.Wrap(err, "slash fraction batch tx")
	}
	if err := validateSlashFractionEthereumSignature(p.SlashFractionEthereumSignature); err != nil {
		return sdkerrors.Wrap(err, "slash fraction ethereum signature")
	}
	if err := validateSlashFractionConflictingEthereumSignature(p.SlashFractionConflictingEthereumSignature); err != nil {
		return sdkerrors.Wrap(err, "slash fraction conflicting ethereum signature")
	}
	if err := validateUnbondSlashingSignerSetTxsWindow(p.UnbondSlashingSignerSetTxsWindow); err != nil {
		return sdkerrors.Wrap(err, "unbond slashing signersettx window")
	}

	return nil
}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyGravityID, &p.GravityId, validateGravityID),
		paramtypes.NewParamSetPair(ParamsStoreKeyContractHash, &p.ContractSourceHash, validateContractHash),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractAddress, &p.BridgeEthereumAddress, validateBridgeContractAddress),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractChainID, &p.BridgeChainId, validateBridgeChainID),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedSignerSetTxsWindow, &p.SignedSignerSetTxsWindow, validateSignedSignerSetTxsWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedBatchesWindow, &p.SignedBatchesWindow, validateSignedBatchesWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeyEthereumSignaturesWindow, &p.EthereumSignaturesWindow, validateEthereumSignaturesWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageBlockTime, &p.AverageBlockTime, validateAverageBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreKeyTargetEthTxTimeout, &p.TargetEthTxTimeout, validateTargetEthTxTimeout),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageEthereumBlockTime, &p.AverageEthereumBlockTime, validateAverageEthereumBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionSignerSetTx, &p.SlashFractionSignerSetTx, validateSlashFractionSignerSetTx),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionBatch, &p.SlashFractionBatch, validateSlashFractionBatch),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionEthereumSignature, &p.SlashFractionEthereumSignature, validateSlashFractionEthereumSignature),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionConflictingEthereumSignature, &p.SlashFractionConflictingEthereumSignature, validateSlashFractionConflictingEthereumSignature),
		paramtypes.NewParamSetPair(ParamStoreUnbondSlashingSignerSetTxsWindow, &p.UnbondSlashingSignerSetTxsWindow, validateUnbondSlashingSignerSetTxsWindow),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	pb, err := p.Marshal()
	if err != nil {
		panic(err)
	}
	p2b, err := p2.Marshal()
	if err != nil {
		panic(err)
	}
	return bytes.Equal(pb, p2b)
}

func validateGravityID(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if _, err := strToFixByteArray(v); err != nil {
		return err
	}
	return nil
}

func validateContractHash(i interface{}) error {
	// TODO: should we validate that the input here is a properly formatted
	// SHA256 (or other) hash?
	if _, ok := i.(string); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBridgeChainID(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateTargetEthTxTimeout(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 60000 {
		return fmt.Errorf("invalid target batch timeout, less than 60 seconds is too short")
	}
	return nil
}

func validateAverageBlockTime(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 100 {
		return fmt.Errorf("invalid average Cosmos block time, too short for latency limitations")
	}
	return nil
}

func validateAverageEthereumBlockTime(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 100 {
		return fmt.Errorf("invalid average Ethereum block time, too short for latency limitations")
	}
	return nil
}

func validateBridgeContractAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !common.IsHexAddress(v) {
		return fmt.Errorf("not an ethereum address: %s", v)
	}
	return nil
}

func validateSignedSignerSetTxsWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateUnbondSlashingSignerSetTxsWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionSignerSetTx(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSignedBatchesWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateEthereumSignaturesWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionBatch(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionEthereumSignature(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionConflictingEthereumSignature(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func strToFixByteArray(s string) ([32]byte, error) {
	var out [32]byte
	if len([]byte(s)) > 32 {
		return out, fmt.Errorf("string too long")
	}
	copy(out[:], s)
	return out, nil
}

func byteArrayToFixByteArray(b []byte) (out [32]byte, err error) {
	if len(b) > 32 {
		return out, fmt.Errorf("array too long")
	}
	copy(out[:], b)
	return out, nil
}
