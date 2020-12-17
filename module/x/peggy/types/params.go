package types

import (
	"bytes"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyPeggyID              = []byte("PeggyID")
	KeyProxyContractHash    = []byte("ProxyContractHash")
	KeyProxyContractAddress = []byte("ProxyContractAddress")
	KeyLogicContractHash    = []byte("LogicContractHash")
	KeyLogicContractAddress = []byte("LogicContractAddress")
	KeyVersion              = []byte("Version")
	KeyStartThreshold       = []byte("StartThreshold")
	KeyBridgeChainID        = []byte("BridgeChainID")
	BootstrapValsetNonce    = []byte("BootstrapValsetNonce")
	KeyBatchInterval        = []byte("BatchInterval")
	KeyBatchNum             = []byte("BatchNum")
	KeyValsetInterval       = []byte("ValsetInterval")
	KeyValsetChange         = []byte("ValsetChange")
)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPeggyID, &p.PeggyId, validatePeggyID),
		paramtypes.NewParamSetPair(KeyProxyContractHash, &p.ProxyContractHash, validateProxyContractHash),
		paramtypes.NewParamSetPair(KeyProxyContractAddress, &p.ProxyContractAddress, validateProxyContractAddress),
		paramtypes.NewParamSetPair(KeyLogicContractHash, &p.LogicContractHash, validateLogicContractHash),
		paramtypes.NewParamSetPair(KeyLogicContractAddress, &p.LogicContractAddress, validateLogicContractAddress),
		paramtypes.NewParamSetPair(KeyVersion, &p.Version, validateVersion),
		paramtypes.NewParamSetPair(KeyStartThreshold, &p.StartThreshold, validateStartThreshold),
		paramtypes.NewParamSetPair(KeyBridgeChainID, &p.BridgeChainId, validateBridgeChainID),
		paramtypes.NewParamSetPair(BootstrapValsetNonce, &p.BootstrapValsetNonce, validateBootstrapValsetNonce),
		paramtypes.NewParamSetPair(KeyBatchInterval, &p.BatchInterval, validateBatchTime),
		paramtypes.NewParamSetPair(KeyBatchNum, &p.BatchNum, validateBatchNum),
		paramtypes.NewParamSetPair(KeyValsetInterval, &p.ValsetInterval, validateUpdateValsetTime),
		paramtypes.NewParamSetPair(KeyValsetChange, &p.ValsetChange, validateUpdateValsetChange),
	}
}

func NewParams(
	peggyID string,
	proxyContractHash string,
	proxyContractAddress string,
	logicContractHash string,
	logicContractAddress string,
	version string,
	startThreshold uint64,
	bridgeChainID uint64,
	bootstrapValsetNonce uint64,
	batchInterval uint64,
	batchNum uint64,
	valsetInterval uint64,
	valsetChange uint64,
) Params {
	return Params{
		PeggyId:              peggyID,
		ProxyContractHash:    proxyContractHash,
		ProxyContractAddress: proxyContractAddress,
		LogicContractHash:    logicContractHash,
		LogicContractAddress: logicContractAddress,
		Version:              version,
		StartThreshold:       startThreshold,
		BridgeChainId:        bridgeChainID,
		BootstrapValsetNonce: bootstrapValsetNonce,
		BatchInterval:        batchInterval,
		BatchNum:             batchNum,
		ValsetInterval:       valsetInterval,
		ValsetChange:         valsetChange,
	}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() *Params {
	return &Params{
		PeggyId:              "",
		ProxyContractHash:    "",
		ProxyContractAddress: "",
		LogicContractHash:    "",
		LogicContractAddress: "",
		Version:              "",
		StartThreshold:       uint64(0),
		BridgeChainId:        uint64(0),
		BootstrapValsetNonce: uint64(0),
		BatchInterval:        uint64(0),
		BatchNum:             uint64(0),
		ValsetInterval:       uint64(0),
		ValsetChange:         uint64(0),
	}
}

func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

func (p Params) ValidateBasic() error {
	// TODO
	return nil
}

func validatePeggyID(i interface{}) error {
	// TODO
	return nil
}

func validateProxyContractHash(i interface{}) error {
	// TODO
	return nil
}

func validateProxyContractAddress(i interface{}) error {
	// TODO
	return nil
}

func validateLogicContractHash(i interface{}) error {
	// TODO
	return nil
}

func validateLogicContractAddress(i interface{}) error {
	// TODO
	return nil
}

func validateVersion(i interface{}) error {
	// TODO
	return nil
}

func validateStartThreshold(i interface{}) error {
	// TODO
	return nil
}

func validateBridgeChainID(i interface{}) error {
	// TODO
	return nil
}

func validateBootstrapValsetNonce(i interface{}) error {
	// TODO
	return nil
}

func validateBatchTime(i interface{}) error {
	// TODO
	return nil
}

func validateBatchNum(i interface{}) error {
	// TODO
	return nil
}

func validateUpdateValsetTime(i interface{}) error {
	// TODO
	return nil
}

func validateUpdateValsetChange(i interface{}) error {
	// TODO
	return nil
}
