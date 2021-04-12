package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gravity"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName
)

var (
	// EthAddressKey indexes cosmos validator account addresses
	// i.e. cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn
	EthAddressKey = []byte{0x1}

	// ValsetRequestKey indexes valset requests by nonce
	ValsetRequestKey = []byte{0x2}

	// ValsetConfirmKey indexes valset confirmations by nonce and the validator account address
	// i.e cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn
	ValsetConfirmKey = []byte{0x3}

	// OracleClaimKey Claim details by nonce and validator address
	// i.e. cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn
	// A claim is named more intuitively than an Attestation, it is literally
	// a validator making a claim to have seen something happen. Claims are
	// attached to attestations which can be thought of as 'the event' that
	// will eventually be executed.
	OracleClaimKey = []byte{0x4}

	// OracleAttestationKey attestation details by nonce and validator address
	// i.e. cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	OracleAttestationKey = []byte{0x5}

	// TransferTxPoolKey indexes the last nonce for the outgoing tx pool
	TransferTxPoolKey = []byte{0x6}

	// SecondIndexTransferTxFeeKey indexes fee amounts by token contract address
	SecondIndexTransferTxFeeKey = []byte{0x9}

	// BatchTxKey indexes outgoing tx batches under a nonce and token address
	BatchTxKey = []byte{0xa}

	// BatchTxBlockKey indexes outgoing tx batches under a block height and token address
	BatchTxBlockKey = []byte{0xb}

	// BatchConfirmKey indexes validator confirmations by token contract address
	BatchConfirmKey = []byte{0xe1}

	// SecondIndexNonceByClaimKey indexes latest nonce for a given claim type
	SecondIndexNonceByClaimKey = []byte{0xf}

	// LastEventNonceByValidatorKey indexes lateset event nonce by validator
	LastEventNonceByValidatorKey = []byte{0xf1}

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0xf2}

	// SequenceKeyPrefix indexes different txids
	SequenceKeyPrefix = []byte{0x7}

	// KeyLastTxPoolID indexes the lastTxPoolID
	KeyLastTxPoolID = append(SequenceKeyPrefix, []byte("lastTxPoolId")...)

	// KeyLastOutgoingBatchID indexes the lastBatchID
	KeyLastOutgoingBatchID = append(SequenceKeyPrefix, []byte("lastBatchId")...)

	// KeyOrchestratorAddress indexes the validator keys for an orchestrator
	KeyOrchestratorAddress = []byte{0xe8}

	// KeyOutgoingLogicCall indexes the outgoing logic calls
	KeyOutgoingLogicCall = []byte{0xde}

	// KeyOutgoingLogicConfirm indexes the outgoing logic confirms
	KeyOutgoingLogicConfirm = []byte{0xae}

	// LastObservedEthereumBlockHeightKey indexes the latest Ethereum block height
	LastObservedEthereumBlockHeightKey = []byte{0xf9}

	// DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key = []byte{0xf3}

	// ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey = []byte{0xf4}

	// LastSlashedValsetNonce indexes the latest slashed valset nonce
	LastSlashedValsetNonce = []byte{0xf5}

	// LatestValsetNonce indexes the latest valset nonce
	LatestValsetNonce = []byte{0xf6}

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0xf7}

	// LastUnBondingBlockHeight indexes the last validator unbonding block height
	LastUnBondingBlockHeight = []byte{0xf8}
)

// GetOrchestratorAddressKey returns the following key format
// prefix
// [0xe8][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetOrchestratorAddressKey(orc sdk.AccAddress) []byte {
	return append(KeyOrchestratorAddress, orc.Bytes()...)
}

// GetEthAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetEthAddressKey(validator sdk.ValAddress) []byte {
	return append(EthAddressKey, validator.Bytes()...)
}

// GetValsetKey returns the following key format
// prefix    nonce
// [0x0][0 0 0 0 0 0 0 1]
func GetValsetKey(nonce uint64) []byte {
	return append(ValsetRequestKey, sdk.Uint64ToBigEndian(nonce)...)
}

// GetValsetConfirmKey returns the following key format
// prefix   nonce                    validator-address
// [0x0][0 0 0 0 0 0 0 1][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
// MARK finish-batches: this is where the key is created in the old (presumed working) code
func GetValsetConfirmKey(nonce uint64, validator sdk.AccAddress) []byte {
	return append(ValsetConfirmKey, append(sdk.Uint64ToBigEndian(nonce), validator.Bytes()...)...)
}

// GetClaimKey returns the following key format
// prefix type               cosmos-validator-address                       nonce                             attestation-details-hash
// [0x0][0 0 0 1][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// The Claim hash identifies a unique event, for example it would have a event nonce, a sender and a receiver. Or an event nonce and a batch nonce. But
// the Claim is stored indexed with the claimer key to make sure that it is unique.
func GetClaimKey(details EthereumEvent) []byte {
	var detailsHash []byte
	if details != nil {
		detailsHash = details.ClaimHash()
	} else {
		panic("No claim without details!")
	}
	claimTypeLen := len([]byte{byte(details.GetType())})
	nonceBz := sdk.Uint64ToBigEndian(details.GetEventNonce())
	key := make([]byte, len(OracleClaimKey)+claimTypeLen+sdk.AddrLen+len(nonceBz)+len(detailsHash))
	copy(key[0:], OracleClaimKey)
	copy(key[len(OracleClaimKey):], []byte{byte(details.GetType())})
	// TODO this is the delegate address, should be stored by the valaddress
	copy(key[len(OracleClaimKey)+claimTypeLen:], details.GetClaimer())
	copy(key[len(OracleClaimKey)+claimTypeLen+sdk.AddrLen:], nonceBz)
	copy(key[len(OracleClaimKey)+claimTypeLen+sdk.AddrLen+len(nonceBz):], detailsHash)
	return key
}

// GetAttestationKey returns the following key format
// prefix     nonce                             claim-details-hash
// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKey(eventNonce uint64, claimHash []byte) []byte {
	key := make([]byte, len(OracleAttestationKey)+len(sdk.Uint64ToBigEndian(0))+len(claimHash))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], sdk.Uint64ToBigEndian(eventNonce))
	copy(key[len(OracleAttestationKey)+len(sdk.Uint64ToBigEndian(0)):], claimHash)
	return key
}

// GetAttestationKeyWithHash returns the following key format
// prefix     nonce                             claim-details-hash
// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKeyWithHash(eventNonce uint64, claimHash []byte) []byte {
	key := make([]byte, len(OracleAttestationKey)+len(sdk.Uint64ToBigEndian(0))+len(claimHash))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], sdk.Uint64ToBigEndian(eventNonce))
	copy(key[len(OracleAttestationKey)+len(sdk.Uint64ToBigEndian(0)):], claimHash)
	return key
}

// GetTransferTxPoolKey returns the following key format
// prefix     id
// [0x6][0 0 0 0 0 0 0 1]
func GetTransferTxPoolKey(id string) []byte {
	// TODO: use hex hash .Bytes() instead ?
	return append(TransferTxPoolKey, []byte(id)...)
}

// GetBatchTxKey returns the following key format
// prefix     nonce                     eth-contract-address
// [0xa][0 0 0 0 0 0 0 1][0xc783df8a850f42e7F7e57013759C285caa701eB6]
func GetBatchTxKey(tokenContract string, nonce uint64) []byte {
	return append(append(BatchTxKey, []byte(tokenContract)...), sdk.Uint64ToBigEndian(nonce)...)
}

// GetBatchTxBlockKey returns the following key format
// prefix     blockheight
// [0xb][0 0 0 0 2 1 4 3]
func GetBatchTxBlockKey(block uint64) []byte {
	return append(BatchTxBlockKey, sdk.Uint64ToBigEndian(block)...)
}

// GetBatchConfirmKey returns the following key format
// prefix           eth-contract-address                BatchNonce                       Validator-address
// [0xe1][0xc783df8a850f42e7F7e57013759C285caa701eB6][0 0 0 0 0 0 0 1][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
// TODO this should be a sdk.ValAddress
func GetBatchConfirmKey(tokenContract string, batchNonce uint64, validator sdk.AccAddress) []byte {
	a := append(sdk.Uint64ToBigEndian(batchNonce), validator.Bytes()...)
	b := append([]byte(tokenContract), a...)
	c := append(BatchConfirmKey, b...)
	return c
}

// GetLastEventNonceByValidatorKey indexes lateset event nonce by validator
// GetLastEventNonceByValidatorKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetLastEventNonceByValidatorKey(validator sdk.ValAddress) []byte {
	return append(LastEventNonceByValidatorKey, validator.Bytes()...)
}

func GetDenomToERC20Key(denom string) []byte {
	return append(DenomToERC20Key, []byte(denom)...)
}

func GetERC20ToDenomKey(erc20 string) []byte {
	return append(ERC20ToDenomKey, []byte(erc20)...)
}

func GetOutgoingLogicCallKey(invalidationId []byte, invalidationNonce uint64) []byte {
	a := append(KeyOutgoingLogicCall, invalidationId...)
	return append(a, sdk.Uint64ToBigEndian(invalidationNonce)...)
}

func GetLogicConfirmKey(invalidationId []byte, invalidationNonce uint64, validator sdk.AccAddress) []byte {
	interm := append(KeyOutgoingLogicConfirm, invalidationId...)
	interm = append(interm, sdk.Uint64ToBigEndian(invalidationNonce)...)
	return append(interm, validator.Bytes()...)
}
