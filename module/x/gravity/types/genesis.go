package types

import (
	"bytes"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParamspace defines the default auth module parameter subspace
const (
	// todo: implement oracle constants as params
	DefaultParamspace = ModuleName
)

var (
	EthereumEventVoteRecordPowerThreshold                          = sdk.NewInt(66)
	ParamsStoreKeyGravityID                                        = []byte("GravityID")
	ParamsStoreKeyContractHash                                     = []byte("ContractHash")
	ParamsStoreKeyStartThreshold                                   = []byte("StartThreshold")
	ParamsStoreKeyBridgeContractAddress                            = []byte("BridgeContractAddress")
	ParamsStoreKeyBridgeContractChainID                            = []byte("BridgeChainID")
	ParamsStoreKeySignedSignerSetTxsWindow                         = []byte("SignedSignerSetTxsWindow")
	ParamsStoreKeySignedSignedBatchTxsWindow                       = []byte("SignedBatchesWindow")
	ParamsStoreKeyTargetBatchTimeout                               = []byte("TargetBatchTimeout")
	ParamsStoreKeyAverageBlockTime                                 = []byte("AverageBlockTime")
	ParamsStoreKeyAverageEthereumBlockTime                         = []byte("AverageEthereumBlockTime")
	ParamsStoreKeySlashFractionSignerSetTx                         = []byte("SlashFractionSignerSetTx")
	ParamsStoreKeySlashFractionBatchTx                             = []byte("SlashFractionBatch")
	ParamStoreKeySlashingSignerSetUnbondWindow                     = []byte("UnbondSlashingSignerSetTxsWindow")
	ParamStoreKeySlashFractionBadEthSignature                      = []byte("SlashFractionBadEthSignature")
	_                                          paramtypes.ParamSet = &Params{}
)

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
		GravityId:                     "defaultgravityid",
		SignedSignerSetTxsWindow:      10000,
		SignedBatchTxsWindow:          10000,
		TargetBatchTimeout:            43200000,
		AverageBlockTime:              5000,
		AverageEthereumBlockTime:      15000,
		SlashFractionSignerSetTx:      sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionBatchTx:          sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionBadEthSignature:  sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashingSignerSetUnbondWindow: 10000,
	}
}

// ValidateBasic checks that the parameters have valid values.
func (p Params) ValidateBasic() error {
	if err := validateGravityID(p.GravityId); err != nil {
		return sdkerrors.Wrap(err, "gravity id")
	}
	if err := validateContractSourceHash(p.ContractSourceHash); err != nil {
		return sdkerrors.Wrap(err, "contract hash")
	}
	if err := validateBridgeEthereumAddress(p.BridgeEthereumAddress); err != nil {
		return sdkerrors.Wrap(err, "bridge contract address")
	}
	if err := validateBridgeChainID(p.BridgeChainId); err != nil {
		return sdkerrors.Wrap(err, "bridge chain id")
	}
	if err := validateTargetBatchTimeout(p.TargetBatchTimeout); err != nil {
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
	if err := validateSignedBatchTxsWindow(p.SignedBatchTxsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSlashFractionSignerSetTx(p.SlashFractionSignerSetTx); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionBatchTx(p.SlashFractionBatchTx); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionBadEthSignature(p.SlashFractionBadEthSignature); err != nil {
		return sdkerrors.Wrap(err, "slash fraction BadEthSignature")
	}
	if err := validateUnbondSlashingSignerSetUnbondWindow(p.SlashingSignerSetUnbondWindow); err != nil {
		return sdkerrors.Wrap(err, "unbond Slashing valset window")
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
		paramtypes.NewParamSetPair(ParamsStoreKeyContractHash, &p.ContractSourceHash, validateContractSourceHash),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractAddress, &p.BridgeEthereumAddress, validateBridgeEthereumAddress),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractChainID, &p.BridgeChainId, validateBridgeChainID),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedSignerSetTxsWindow, &p.SignedSignerSetTxsWindow, validateSignedSignerSetTxsWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedSignedBatchTxsWindow, &p.SignedBatchTxsWindow, validateSignedBatchTxsWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageBlockTime, &p.AverageBlockTime, validateAverageBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreKeyTargetBatchTimeout, &p.TargetBatchTimeout, validateTargetBatchTimeout),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageEthereumBlockTime, &p.AverageEthereumBlockTime, validateAverageEthereumBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreKeySlashFractionSignerSetTx, &p.SlashFractionSignerSetTx, validateSlashFractionSignerSetTx),
		paramtypes.NewParamSetPair(ParamsStoreKeySlashFractionBatchTx, &p.SlashFractionBatchTx, validateSlashFractionBatchTx),
		paramtypes.NewParamSetPair(ParamStoreKeySlashFractionBadEthSignature, &p.SlashFractionBadEthSignature, validateSlashFractionBadEthSignature),
		paramtypes.NewParamSetPair(ParamStoreKeySlashingSignerSetUnbondWindow, &p.SlashingSignerSetUnbondWindow, validateUnbondSlashingSignerSetUnbondWindow),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
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

func validateContractSourceHash(i interface{}) error {
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

func validateTargetBatchTimeout(i interface{}) error {
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

func validateBridgeEthereumAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := ValidateEthAddress(v); err != nil {
		// TODO: ensure that empty addresses are valid in params
		if !strings.Contains(err.Error(), "empty") {
			return err
		}
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

func validateUnbondSlashingSignerSetUnbondWindow(i interface{}) error {
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

func validateSignedBatchTxsWindow(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionBatchTx(i interface{}) error {
	// TODO: do we want to set some bounds on this value?
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionBadEthSignature(i interface{}) error {
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
