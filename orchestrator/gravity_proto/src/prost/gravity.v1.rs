/// Attestation is an event that is pending of confirmation by 2/3 of the signer set.
/// The event is then attested and executed vy the state machine once the required
/// threshold is met.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Attestation {
    /// event unique identifier
    #[prost(bytes="vec", tag="1")]
    pub event_id: ::prost::alloc::vec::Vec<u8>,
    /// set of the validator operators address in bech32 format that attest in
    /// favor of this event.
    #[prost(string, repeated, tag="2")]
    pub votes: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    /// amount of voting power in support of this event execution
    #[prost(int64, tag="3")]
    pub attested_power: i64,
    /// height at which the event was attested an executed
    #[prost(uint64, tag="4")]
    pub height: u64,
}
/// DepositEvent is submitted when more than 66% of the active
/// Cosmos validator set has claimed to have seen a deposit
/// on Ethereum. ERC20 coins are minted to the receiver address
/// address.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DepositEvent {
    /// event nonce for replay protection
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    /// ethereum hex address of the contract
    #[prost(string, tag="2")]
    pub token_contract: ::prost::alloc::string::String,
    /// amount of tokens deposited on Ethereum
    #[prost(string, tag="3")]
    pub amount: ::prost::alloc::string::String,
    /// ethereum sender address in hex format
    #[prost(string, tag="4")]
    pub ethereum_sender: ::prost::alloc::string::String,
    /// cosmos bech32 account address of the receiver
    #[prost(string, tag="5")]
    pub cosmos_receiver: ::prost::alloc::string::String,
    /// etherereum block height at which the event was observed
    #[prost(uint64, tag="6")]
    pub ethereum_height: u64,
}
/// WithdrawEvent claims that a batch of withdrawal
/// operations on the bridge contract was executed.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct WithdrawEvent {
    /// transaction identifier of the batch tx
    #[prost(bytes="vec", tag="1")]
    pub tx_id: ::prost::alloc::vec::Vec<u8>,
    /// event nonce of the batch tx on Cosmos
    #[prost(uint64, tag="2")]
    pub nonce: u64,
    /// ethereum hex address of the contract
    #[prost(string, tag="3")]
    pub token_contract: ::prost::alloc::string::String,
    /// etherereum block height at which the event was observed
    #[prost(uint64, tag="4")]
    pub ethereum_height: u64,
}
/// LogicCallExecutedEvent describes a logic call that has been
/// successfully executed on Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LogicCallExecutedEvent {
    /// event nonce for replay protection
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    /// invalidation ID allows to withdraw funds in case the logic call fails on the
    /// ethereum side
    #[prost(bytes="vec", tag="2")]
    pub invalidation_id: ::prost::alloc::vec::Vec<u8>,
    /// TODO: explain
    #[prost(uint64, tag="3")]
    pub invalidation_nonce: u64,
    /// etherereum block height at which the event was observed
    #[prost(uint64, tag="4")]
    pub ethereum_height: u64,
}
/// CosmosERC20DeployedEvent is submitted when an ERC20 contract
/// for a Cosmos SDK coin has been deployed on Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CosmosErc20DeployedEvent {
    /// event nonce for replay protection
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    /// cosmos SDK coin denomination
    #[prost(string, tag="2")]
    pub cosmos_denom: ::prost::alloc::string::String,
    /// ethereum ERC20 contract address in hex format
    #[prost(string, tag="3")]
    pub token_contract: ::prost::alloc::string::String,
    /// name of the token
    #[prost(string, tag="4")]
    pub name: ::prost::alloc::string::String,
    /// symbol or tick of the token
    #[prost(string, tag="5")]
    pub symbol: ::prost::alloc::string::String,
    /// number of decimals the token supports (i.e precision)
    #[prost(uint64, tag="6")]
    pub decimals: u64,
    /// etherereum block height at which the event was observed
    #[prost(uint64, tag="7")]
    pub ethereum_height: u64,
}
/// EthereumInfo defines the latest observed ethereum block height and the
/// corresponding timestamp value in nanoseconds.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EthereumInfo {
    /// timestamp in nanoseconds
    #[prost(message, optional, tag="1")]
    pub timestamp: ::core::option::Option<::prost_types::Timestamp>,
    /// ethereum block height
    #[prost(uint64, tag="2")]
    pub height: u64,
}
/// EthSigner represents a cosmos validator with its corresponding bridge operator
/// ethereum address and its staking consensus power.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EthSigner {
    /// staking consensus power of the validator
    #[prost(int64, tag="1")]
    pub power: i64,
    /// bridge operator ethereum address in hex format
    #[prost(string, tag="2")]
    pub ethereum_address: ::prost::alloc::string::String,
}
/// EthSignerSet is the Ethereum Bridge multisig set that relays transactions
/// the two chains. The staking validators keep ethereum keys which are used to
/// check signatures on Ethereum in order to get significant gas savings.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EthSignerSet {
    /// set of signers. Sorted by power
    #[prost(message, repeated, tag="1")]
    pub signers: ::prost::alloc::vec::Vec<EthSigner>,
    /// TODO: which height? cosmos? This should be the key
    #[prost(uint64, tag="2")]
    pub height: u64,
}
/// BatchTx represents a batch of transactions going from Cosmos to Ethereum
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTx {
    /// tx nonce for replay protection
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    /// timeout in seconds  // TODO: double check
    #[prost(uint64, tag="2")]
    pub timeout: u64,
    /// transaction identifiers of the transfer txs included in this batch
    #[prost(bytes="vec", repeated, tag="3")]
    pub transactions: ::prost::alloc::vec::Vec<::prost::alloc::vec::Vec<u8>>,
    /// ethereum contract address of the batch contract in hex format
    #[prost(string, tag="4")]
    pub token_contract: ::prost::alloc::string::String,
    /// ethereum block height // TODO: double check
    #[prost(uint64, tag="5")]
    pub block: u64,
}
/// TransferTx represents an individual transfer from Cosmos to Ethereum
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TransferTx {
    /// tx nonce for replay protection
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    /// cosmos account address of the sender in bech32 format
    #[prost(string, tag="2")]
    pub sender: ::prost::alloc::string::String,
    /// ethereum recipient address in hex format
    #[prost(string, tag="3")]
    pub ethereum_recipient: ::prost::alloc::string::String,
    /// amount of the transfer represented as an sdk.Coin. The coin denomination
    /// must correspond to a valid ERC20 token contract address
    #[prost(message, optional, tag="4")]
    pub erc20_token: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    /// transfer fees for the bridge orchestrators, represented as an sdk.Coin.
    /// The coin denomination must correspond to a valid ERC20 token contract address
    #[prost(message, optional, tag="5")]
    pub erc20_fee: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
/// LogicCallTx represents an individual arbitratry logic call transaction from
/// Cosmos to Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LogicCallTx {
    /// tx nonce for replay protection
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    /// erc20 tokens represented as sdk.Coins
    #[prost(message, repeated, tag="2")]
    pub tokens: ::prost::alloc::vec::Vec<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    /// erc20 tokens represented as sdk.Coins used as fees for the bridge orchestrators.
    #[prost(message, repeated, tag="3")]
    pub fees: ::prost::alloc::vec::Vec<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    /// ethereum contract address of the arbitrary logic contract in hex format
    #[prost(string, tag="4")]
    pub logic_contract_address: ::prost::alloc::string::String,
    /// ABI payload of the smart contract function call
    #[prost(bytes="vec", tag="5")]
    pub payload: ::prost::alloc::vec::Vec<u8>,
    /// timeout in seconds  // TODO: double check
    #[prost(uint64, tag="6")]
    pub timeout: u64,
}
/// TransactionIDs defines a protobuf message for storing transfer tx ids.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TransactionIDs {
    /// slice of transfer transaction identifiers
    #[prost(bytes="vec", repeated, tag="1")]
    pub ids: ::prost::alloc::vec::Vec<::prost::alloc::vec::Vec<u8>>,
}
/// ConfirmLogicCall ...
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ConfirmLogicCall {
    #[prost(bytes="vec", tag="1")]
    pub invalidation_id: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="2")]
    pub invalidation_nonce: u64,
    #[prost(string, tag="3")]
    pub eth_signer: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub orchestrator_address: ::prost::alloc::string::String,
    #[prost(bytes="vec", tag="5")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
}
/// ConfirmBatch an orchestrator confirms a batch transaction by signing
/// with the ethereum keys on the signer set.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ConfirmBatch {
    #[prost(string, tag="1")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(uint64, tag="2")]
    pub nonce: u64,
    #[prost(string, tag="3")]
    pub eth_signer: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub orchestrator_address: ::prost::alloc::string::String,
    #[prost(bytes="vec", tag="5")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
}
/// ConfirmSignerSet submits a signature of the validator set at a given block height. A validator
/// must first call MsgSetEthAddress to set their Ethereum address to be used for signing.
/// Finally validators sign the
/// validator set, powers, and Ethereum addresses of the entire validator set at the height of a
/// ValsetRequest and submit that signature with this message.
///
/// If a sufficient number of validators (66% of voting power) (A) have set Ethereum addresses and
/// (B) submit ValsetConfirm messages with their signatures it is then possible for anyone to view
/// these signatures in the chain store and submit them to Ethereum to update the validator set
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ConfirmSignerSet {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(string, tag="2")]
    pub eth_signer: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub orchestrator_address: ::prost::alloc::string::String,
    #[prost(bytes="vec", tag="4")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
}
/// MsgTransfer submits a transfer attempt to bridge an asset over to Ethereum.
/// The transfer will be stored and then included in a batch and then
/// submitted to Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgTransfer {
    /// cosmos account address of the sender in bech32 format
    #[prost(string, tag="1")]
    pub sender: ::prost::alloc::string::String,
    /// ethereum hex address of the recipient
    #[prost(string, tag="2")]
    pub eth_recipient: ::prost::alloc::string::String,
    /// the SDK coin to send across the bridge to Ethereum. This can be either an
    /// ERC20 token voucher or a native cosmos denomination (including IBC vouchers).
    #[prost(message, optional, tag="3")]
    pub amount: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    /// the fee paid for the transaction relay accross the bridge to Ethereum.
    /// NOTE: this is distinct from the Cosmos transaction fee paid, so a successful
    /// transfer has two layers of fees for the user (Cosmos & Bridge).
    /// TODO: specify if this needs to be an ERC20 or not.
    #[prost(message, optional, tag="4")]
    pub bridge_fee: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
/// MsgTransferResponse returns the transfer transaction ID which will be included
/// in the batch tx.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgTransferResponse {
    /// transfer tx unique identifier
    #[prost(bytes="vec", tag="1")]
    pub tx_id: ::prost::alloc::vec::Vec<u8>,
}
/// MsgCancelTransfer allows the sender to cancel its own outgoing transfer tx
/// and recieve a refund of the tokens and bridge fees. This tx will only succeed
/// if the transfer tx hasn't been batched to be processed and relayed to Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgCancelTransfer {
    /// transfer tx unique identifier
    #[prost(bytes="vec", tag="1")]
    pub tx_id: ::prost::alloc::vec::Vec<u8>,
    /// cosmos account address of the sender in bech32 format
    #[prost(string, tag="2")]
    pub sender: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgCancelTransferResponse {
}
/// MsgRequestBatch requests a batch of transactions with a given coin denomination to send across
/// the bridge to Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRequestBatch {
    /// cosmos account address of the orchestrator in bech32 format
    #[prost(string, tag="1")]
    pub orchestrator_address: ::prost::alloc::string::String,
    /// coin denomination
    #[prost(string, tag="2")]
    pub denom: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRequestBatchResponse {
}
/// MsgSubmitConfirm
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitConfirm {
    /// confirmation for batch txs, logic call txs or signer sets
    #[prost(message, optional, tag="1")]
    pub confirm: ::core::option::Option<::prost_types::Any>,
    /// cosmos account address of the orchestrator signer in bech32 format
    #[prost(string, tag="2")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitConfirmResponse {
}
/// MsgSubmitEvent
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEvent {
    /// event requested observed by a single validator orchestrator on Ethereum,
    /// which will then need to be
    #[prost(message, optional, tag="1")]
    pub event: ::core::option::Option<::prost_types::Any>,
    /// cosmos account address of the orchestrator signer in bech32 format
    #[prost(string, tag="2")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEventResponse {
}
/// MsgDelegateKey allows validators to delegate their voting responsibilities
/// to a given orchestrator address. This key is then used as an optional
/// authentication method for attesting events from Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgDelegateKey {
    /// validator operator address in bech32 format
    #[prost(string, tag="1")]
    pub validator_address: ::prost::alloc::string::String,
    /// cosmos account address of the orchestrator in bech32 format that
    /// references the key that is being delegated to
    #[prost(string, tag="2")]
    pub orchestrator_address: ::prost::alloc::string::String,
    /// ethereum hex address of the used by the orchestrator
    #[prost(string, tag="3")]
    pub eth_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgDelegateKeyResponse {
}
# [doc = r" Generated client implementations."] pub mod msg_client { # ! [allow (unused_variables , dead_code , missing_docs)] use tonic :: codegen :: * ; # [doc = " Msg defines the state transitions possible within gravity"] pub struct MsgClient < T > { inner : tonic :: client :: Grpc < T > , } impl MsgClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > MsgClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + HttpBody + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as HttpBody > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor (inner : T , interceptor : impl Into < tonic :: Interceptor >) -> Self { let inner = tonic :: client :: Grpc :: with_interceptor (inner , interceptor) ; Self { inner } } pub async fn transfer (& mut self , request : impl tonic :: IntoRequest < super :: MsgTransfer > ,) -> Result < tonic :: Response < super :: MsgTransferResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/Transfer") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn cancel_transfer (& mut self , request : impl tonic :: IntoRequest < super :: MsgCancelTransfer > ,) -> Result < tonic :: Response < super :: MsgCancelTransferResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/CancelTransfer") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn request_batch (& mut self , request : impl tonic :: IntoRequest < super :: MsgRequestBatch > ,) -> Result < tonic :: Response < super :: MsgRequestBatchResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/RequestBatch") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn submit_confirm (& mut self , request : impl tonic :: IntoRequest < super :: MsgSubmitConfirm > ,) -> Result < tonic :: Response < super :: MsgSubmitConfirmResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/SubmitConfirm") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn submit_event (& mut self , request : impl tonic :: IntoRequest < super :: MsgSubmitEvent > ,) -> Result < tonic :: Response < super :: MsgSubmitEventResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/SubmitEvent") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn set_delegate_key (& mut self , request : impl tonic :: IntoRequest < super :: MsgDelegateKey > ,) -> Result < tonic :: Response < super :: MsgDelegateKeyResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/SetDelegateKey") ; self . inner . unary (request . into_request () , path , codec) . await } } impl < T : Clone > Clone for MsgClient < T > { fn clone (& self) -> Self { Self { inner : self . inner . clone () , } } } impl < T > std :: fmt :: Debug for MsgClient < T > { fn fmt (& self , f : & mut std :: fmt :: Formatter < '_ >) -> std :: fmt :: Result { write ! (f , "MsgClient {{ ... }}") } } }/// Params represent the Gravity genesis and store parameters
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
    /// address of the bridge contract on the EVM chain
    #[prost(string, tag="1")]
    pub bridge_contract_address: ::prost::alloc::string::String,
    /// unique identifier of the EVM chain
    #[prost(uint64, tag="2")]
    pub bridge_chain_id: u64,
    /// target value for when batch transactions time out on Ethereum
    #[prost(uint64, tag="3")]
    pub target_batch_timeout: u64,
    /// Average Cosmos block time used to compute batch timeout
    #[prost(uint64, tag="4")]
    pub average_block_time: u64,
    /// Average ethereum block time used to compute batch timeout
    #[prost(uint64, tag="5")]
    pub average_ethereum_block_time: u64,
    /// amount of blocks of the rolling window required to submit a signature for a signer set confirmation.
    #[prost(uint64, tag="6")]
    pub signer_set_window: u64,
    /// amount of blocks of the rolling window required to submit a signature for a batch transaction.
    #[prost(uint64, tag="7")]
    pub batch_tx_window: u64,
    /// amount of blocks of the rolling window required to attest an ethereum event.
    #[prost(uint64, tag="8")]
    pub event_window: u64,
    #[prost(uint64, tag="9")]
    pub unbonding_window: u64,
    /// max amount of transactions batched
    #[prost(uint64, tag="10")]
    pub batch_size: u64,
    /// slashing fraction for not signing a signerset confirmation
    #[prost(string, tag="11")]
    pub slash_fraction_signer_set: ::prost::alloc::string::String,
    /// slashing fraction for not signing an outgoing batch transaction to ethereum
    #[prost(string, tag="12")]
    pub slash_fraction_batch: ::prost::alloc::string::String,
    /// slashing fraction for not signing events
    #[prost(string, tag="13")]
    pub slash_fraction_event: ::prost::alloc::string::String,
    /// slashing fraction for submitting a conflicting event from Ethereum
    #[prost(string, tag="14")]
    pub slash_fraction_conflicting_event: ::prost::alloc::string::String,
    /// maximum allowed power difference between the latest and the current ethereum 
    /// signer set before submitting a new signer set update request.
    #[prost(string, tag="15")]
    pub max_signer_set_power_diff: ::prost::alloc::string::String,
}
/// ERC20ToDenom records the relationship between an ERC20 token contract and the
/// denomination of the corresponding Cosmos coin.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20ToDenom {
    /// token contract address in hex format
    #[prost(string, tag="1")]
    pub erc20_address: ::prost::alloc::string::String,
    /// coin denomination
    #[prost(string, tag="2")]
    pub denom: ::prost::alloc::string::String,
}
/// GenesisState struct
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    /// bridge id is a random 32 byte salt value to prevent signature reuse across different
    /// instances of the ethereum bridge. This value IS USED on the bridge smart
    /// contracts.
    ///
    /// TODO: is this like the IBC client ID but for the bridge?
    #[prost(bytes="vec", tag="1")]
    pub bridge_id: ::prost::alloc::vec::Vec<u8>,
    #[prost(message, optional, tag="2")]
    pub params: ::core::option::Option<Params>,
    #[prost(uint64, tag="3")]
    pub last_observed_nonce: u64,
    #[prost(message, repeated, tag="4")]
    pub signer_sets: ::prost::alloc::vec::Vec<EthSignerSet>,
    /// requested batch transactions
    #[prost(message, repeated, tag="5")]
    pub batch_txs: ::prost::alloc::vec::Vec<BatchTx>,
    #[prost(message, repeated, tag="6")]
    pub logic_call_txs: ::prost::alloc::vec::Vec<LogicCallTx>,
    /// unbatched transfer transactions
    ///
    /// TODO: use any for confirms
    #[prost(message, repeated, tag="7")]
    pub transfer_txs: ::prost::alloc::vec::Vec<TransferTx>,
    #[prost(message, repeated, tag="8")]
    pub signer_set_confirms: ::prost::alloc::vec::Vec<ConfirmSignerSet>,
    #[prost(message, repeated, tag="9")]
    pub batch_confirms: ::prost::alloc::vec::Vec<ConfirmBatch>,
    #[prost(message, repeated, tag="10")]
    pub logic_call_confirms: ::prost::alloc::vec::Vec<ConfirmLogicCall>,
    /// TODO: proto.Any ethereum eventss
    #[prost(message, repeated, tag="11")]
    pub attestations: ::prost::alloc::vec::Vec<Attestation>,
    #[prost(message, repeated, tag="12")]
    pub delegate_keys: ::prost::alloc::vec::Vec<MsgDelegateKey>,
    #[prost(message, repeated, tag="13")]
    pub erc20_to_denoms: ::prost::alloc::vec::Vec<Erc20ToDenom>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    #[prost(message, optional, tag="1")]
    pub params: ::core::option::Option<Params>,
}
# [doc = r" Generated client implementations."] pub mod query_client { # ! [allow (unused_variables , dead_code , missing_docs)] use tonic :: codegen :: * ; # [doc = " Query defines the gRPC querier service"] pub struct QueryClient < T > { inner : tonic :: client :: Grpc < T > , } impl QueryClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > QueryClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + HttpBody + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as HttpBody > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor (inner : T , interceptor : impl Into < tonic :: Interceptor >) -> Self { let inner = tonic :: client :: Grpc :: with_interceptor (inner , interceptor) ; Self { inner } } pub async fn params (& mut self , request : impl tonic :: IntoRequest < super :: QueryParamsRequest > ,) -> Result < tonic :: Response < super :: QueryParamsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/Params") ; self . inner . unary (request . into_request () , path , codec) . await } } impl < T : Clone > Clone for QueryClient < T > { fn clone (& self) -> Self { Self { inner : self . inner . clone () , } } } impl < T > std :: fmt :: Debug for QueryClient < T > { fn fmt (& self , f : & mut std :: fmt :: Formatter < '_ >) -> std :: fmt :: Result { write ! (f , "QueryClient {{ ... }}") } } }
