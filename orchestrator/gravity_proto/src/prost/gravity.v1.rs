/// EthereumEventVoteRecord is an event that is pending of confirmation by 2/3 of the signer set.
/// The event is then attested and executed in the state machine once the required
/// threshold is met.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EthereumEventVoteRecord {
    #[prost(message, optional, tag="1")]
    pub event: ::core::option::Option<::prost_types::Any>,
    #[prost(string, repeated, tag="2")]
    pub votes: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    #[prost(bool, tag="3")]
    pub accepted: bool,
}
/// LatestEthereumBlockHeight defines the latest observed ethereum block height and the
/// corresponding timestamp value in nanoseconds.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LatestEthereumBlockHeight {
    #[prost(uint64, tag="1")]
    pub ethereum_height: u64,
    #[prost(uint64, tag="2")]
    pub cosmos_height: u64,
}
/// EthereumSigner represents a cosmos validator with its corresponding bridge operator
/// ethereum address and its staking consensus power.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EthereumSigner {
    #[prost(int64, tag="1")]
    pub power: i64,
    #[prost(string, tag="2")]
    pub ethereum_address: ::prost::alloc::string::String,
}
/// UpdateSignerSetTx is the Ethereum Bridge multisig set that relays transactions
/// the two chains. The staking validators keep ethereum keys which are used to
/// check signatures on Ethereum in order to get significant gas savings.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateSignerSetTx {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(message, repeated, tag="2")]
    pub signers: ::prost::alloc::vec::Vec<EthereumSigner>,
}
/// BatchTx represents a batch of transactions going from Cosmos to Ethereum. Batch
/// txs are are identified by a unique hash and the token contract that is shared by
/// all the SendToEthereum
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTx {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(uint64, tag="2")]
    pub timeout: u64,
    #[prost(bytes="vec", repeated, tag="3")]
    pub transactions: ::prost::alloc::vec::Vec<::prost::alloc::vec::Vec<u8>>,
    #[prost(string, tag="4")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(uint64, tag="5")]
    pub ethereum_block: u64,
}
/// SendToEthereumTx represents an individual SendToEthereum from Cosmos to Ethereum
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SendToEthereumTx {
    #[prost(uint64, tag="1")]
    pub id: u64,
    #[prost(string, tag="2")]
    pub sender: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub ethereum_recipient: ::prost::alloc::string::String,
    #[prost(message, optional, tag="4")]
    pub erc20_token: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, optional, tag="5")]
    pub erc20_fee: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
/// ContractCallTx represents an individual arbitratry logic call transaction from
/// Cosmos to Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTx {
    #[prost(uint64, tag="1")]
    pub invalidation_nonce: u64,
    #[prost(bytes="vec", tag="2")]
    pub invalidation_scope: ::prost::alloc::vec::Vec<u8>,
    #[prost(string, tag="3")]
    pub contract_call_address: ::prost::alloc::string::String,
    #[prost(bytes="vec", tag="4")]
    pub payload: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="5")]
    pub timeout: u64,
    #[prost(message, repeated, tag="6")]
    pub tokens: ::prost::alloc::vec::Vec<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, repeated, tag="7")]
    pub fees: ::prost::alloc::vec::Vec<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
/// MsgSendToEthereum submits a SendToEthereum attempt to bridge an asset over to Ethereum.
/// The SendToEthereum will be stored and then included in a batch and then
/// submitted to Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSendToEthereum {
    #[prost(string, tag="1")]
    pub sender: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub eth_recipient: ::prost::alloc::string::String,
    #[prost(message, optional, tag="3")]
    pub amount: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, optional, tag="4")]
    pub bridge_fee: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
/// MsgSendToEthereumResponse returns the SendToEthereum transaction ID which will be included
/// in the batch tx.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSendToEthereumResponse {
    #[prost(uint64, tag="1")]
    pub id: u64,
}
/// MsgCancelSendToEthereum allows the sender to cancel its own outgoing SendToEthereum tx
/// and recieve a refund of the tokens and bridge fees. This tx will only succeed
/// if the SendToEthereum tx hasn't been batched to be processed and relayed to Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgCancelSendToEthereum {
    #[prost(uint64, tag="1")]
    pub id: u64,
    #[prost(string, tag="2")]
    pub sender: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgCancelSendToEthereumResponse {
}
/// MsgRequestBatch requests a batch of transactions with a given coin denomination to send across
/// the bridge to Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRequestBatchTx {
    #[prost(string, tag="1")]
    pub denom: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRequestBatchTxResponse {
}
/// MsgSubmitEthereumSignature submits an ethereum signature for a given validator
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEthereumSignature {
    /// TODO: can we make this take an array?
    #[prost(message, optional, tag="1")]
    pub signature: ::core::option::Option<::prost_types::Any>,
    #[prost(string, tag="2")]
    pub signer: ::prost::alloc::string::String,
}
/// ContractCallTxSignature is a signature on behalf of a validator for a ContractCallTx.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxSignature {
    #[prost(bytes="vec", tag="1")]
    pub invalidation_id: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="2")]
    pub invalidation_nonce: u64,
    #[prost(string, tag="3")]
    pub eth_signer: ::prost::alloc::string::String,
    #[prost(bytes="vec", tag="4")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
}
/// BatchTxSignature is a signature on behalf of a validator for a BatchTx.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxSignature {
    #[prost(string, tag="1")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(uint64, tag="2")]
    pub nonce: u64,
    #[prost(string, tag="3")]
    pub eth_signer: ::prost::alloc::string::String,
    #[prost(bytes="vec", tag="4")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
}
/// UpdateSignerSetTxSignature is a signature on behalf of a validator for a UpdateSignerSetTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateSignerSetTxSignature {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(string, tag="2")]
    pub eth_signer: ::prost::alloc::string::String,
    #[prost(bytes="vec", tag="3")]
    pub signature: ::prost::alloc::vec::Vec<u8>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEthereumSignatureResponse {
}
/// MsgSubmitEthereumEvent
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEthereumEvent {
    #[prost(message, optional, tag="1")]
    pub event: ::core::option::Option<::prost_types::Any>,
    #[prost(string, tag="2")]
    pub signer: ::prost::alloc::string::String,
}
/// SendToCosmosEvent is submitted when the SendToCosmosEvent is emitted by they gravity contract. 
/// ERC20 representation coins are minted to the cosmosreceiver address.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SendToCosmosEvent {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
    #[prost(string, tag="2")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub amount: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub ethereum_sender: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub cosmos_receiver: ::prost::alloc::string::String,
    #[prost(uint64, tag="6")]
    pub ethereum_height: u64,
}
/// BatchExecutedEvent claims that a batch of BatchTxExecutedal operations on the bridge contract
/// was executed successfully on ETH
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchExecutedEvent {
    #[prost(string, tag="1")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(uint64, tag="2")]
    pub event_nonce: u64,
    #[prost(uint64, tag="3")]
    pub ethereum_height: u64,
}
// ContractCallExecutedEvent describes a contract call that has been
// successfully executed on Ethereum.

/// NOTE: bytes.HexBytes is supposed to "help" with json encoding/decoding investigate?
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallExecutedEvent {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
    #[prost(bytes="vec", tag="2")]
    pub invalidation_id: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="3")]
    pub invalidation_nonce: u64,
    #[prost(uint64, tag="4")]
    pub ethereum_height: u64,
}
/// ERC20DeployedEvent is submitted when an ERC20 contract
/// for a Cosmos SDK coin has been deployed on Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20DeployedEvent {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
    #[prost(string, tag="2")]
    pub cosmos_denom: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub token_contract: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub erc20_name: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub erc20_symbol: ::prost::alloc::string::String,
    #[prost(uint64, tag="6")]
    pub erc20_decimals: u64,
    #[prost(uint64, tag="7")]
    pub ethereum_height: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitEthereumEventResponse {
}
/// MsgDelegateKey allows validators to delegate their voting responsibilities
/// to a given orchestrator address. This key is then used as an optional
/// authentication method for attesting events from Ethereum.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgDelegateKeys {
    #[prost(string, tag="1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub orchestrator_address: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub eth_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgDelegateKeysResponse {
}
# [doc = r" Generated client implementations."] pub mod msg_client { # ! [allow (unused_variables , dead_code , missing_docs)] use tonic :: codegen :: * ; # [doc = " Msg defines the state transitions possible within gravity"] pub struct MsgClient < T > { inner : tonic :: client :: Grpc < T > , } impl MsgClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > MsgClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + HttpBody + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as HttpBody > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor (inner : T , interceptor : impl Into < tonic :: Interceptor >) -> Self { let inner = tonic :: client :: Grpc :: with_interceptor (inner , interceptor) ; Self { inner } } pub async fn send_to_ethereum (& mut self , request : impl tonic :: IntoRequest < super :: MsgSendToEthereum > ,) -> Result < tonic :: Response < super :: MsgSendToEthereumResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/SendToEthereum") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn cancel_send_to_ethereum (& mut self , request : impl tonic :: IntoRequest < super :: MsgCancelSendToEthereum > ,) -> Result < tonic :: Response < super :: MsgCancelSendToEthereumResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/CancelSendToEthereum") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn request_batch_tx (& mut self , request : impl tonic :: IntoRequest < super :: MsgRequestBatchTx > ,) -> Result < tonic :: Response < super :: MsgRequestBatchTxResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/RequestBatchTx") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn submit_ethereum_signature (& mut self , request : impl tonic :: IntoRequest < super :: MsgSubmitEthereumSignature > ,) -> Result < tonic :: Response < super :: MsgSubmitEthereumSignatureResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/SubmitEthereumSignature") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn submit_ethereum_event (& mut self , request : impl tonic :: IntoRequest < super :: MsgSubmitEthereumEvent > ,) -> Result < tonic :: Response < super :: MsgSubmitEthereumEventResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/SubmitEthereumEvent") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn set_delegate_keys (& mut self , request : impl tonic :: IntoRequest < super :: MsgDelegateKeys > ,) -> Result < tonic :: Response < super :: MsgDelegateKeysResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Msg/SetDelegateKeys") ; self . inner . unary (request . into_request () , path , codec) . await } } impl < T : Clone > Clone for MsgClient < T > { fn clone (& self) -> Self { Self { inner : self . inner . clone () , } } } impl < T > std :: fmt :: Debug for MsgClient < T > { fn fmt (& self , f : & mut std :: fmt :: Formatter < '_ >) -> std :: fmt :: Result { write ! (f , "MsgClient {{ ... }}") } } }/// Params represent the Gravity genesis and store parameters
/// gravity_id:
/// a random 32 byte value to prevent signature reuse, for example if the
/// cosmos validators decided to use the same Ethereum keys for another chain
/// also running Gravity we would not want it to be possible to play a deposit
/// from chain A back on chain B's Gravity. This value IS USED ON ETHEREUM so
/// it must be set in your genesis.json before launch and not changed after
/// deploying Gravity
///
/// contract_hash:
/// the code hash of a known good version of the Gravity contract
/// solidity code. This can be used to verify the correct version
/// of the contract has been deployed. This is a reference value for
/// goernance action only it is never read by any Gravity code
///
/// bridge_ethereum_address:
/// is address of the bridge contract on the Ethereum side, this is a
/// reference value for governance only and is not actually used by any
/// Gravity code
///
/// bridge_chain_id:
/// the unique identifier of the Ethereum chain, this is a reference value
/// only and is not actually used by any Gravity code
///
/// These reference values may be used by future Gravity client implemetnations
/// to allow for saftey features or convenience features like the Gravity address
/// in your relayer. A relayer would require a configured Gravity address if
/// governance had not set the address on the chain it was relaying for.
///
/// signed_valsets_window
/// signed_batches_window
/// signed_claims_window
///
/// These values represent the time in blocks that a validator has to submit
/// a signature for a batch or valset, or to submit a claim for a particular
/// attestation nonce. In the case of attestations this clock starts when the
/// attestation is created, but only allows for slashing once the event has passed
///
/// target_batch_timeout:
///
/// This is the 'target' value for when batches time out, this is a target becuase
/// Ethereum is a probabalistic chain and you can't say for sure what the block
/// frequency is ahead of time.
///
/// average_block_time
/// average_ethereum_block_time
///
/// These values are the average Cosmos block time and Ethereum block time repsectively
/// and they are used to copute what the target batch timeout is. It is important that
/// governance updates these in case of any major, prolonged change in the time it takes
/// to produce a block
///
/// slash_fraction_valset
/// slash_fraction_batch
/// slash_fraction_claim
/// slash_fraction_conflicting_claim
///
/// The slashing fractions for the various gravity related slashing conditions. The first three
/// refer to not submitting a particular message, the third for submitting a different claim
/// for the same Ethereum event
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
    #[prost(string, tag="1")]
    pub gravity_id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub contract_source_hash: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub bridge_ethereum_address: ::prost::alloc::string::String,
    #[prost(uint64, tag="5")]
    pub bridge_chain_id: u64,
    #[prost(uint64, tag="6")]
    pub signed_valsets_window: u64,
    #[prost(uint64, tag="7")]
    pub signed_batches_window: u64,
    #[prost(uint64, tag="8")]
    pub signed_claims_window: u64,
    #[prost(uint64, tag="10")]
    pub target_batch_timeout: u64,
    #[prost(uint64, tag="11")]
    pub average_block_time: u64,
    #[prost(uint64, tag="12")]
    pub average_ethereum_block_time: u64,
    #[prost(bytes="vec", tag="13")]
    pub slash_fraction_valset: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes="vec", tag="14")]
    pub slash_fraction_batch: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes="vec", tag="15")]
    pub slash_fraction_claim: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes="vec", tag="16")]
    pub slash_fraction_conflicting_claim: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="17")]
    pub unbond_slashing_valsets_window: u64,
}
/// GenesisState struct
/// TODO: this need to be audited and potentially simplified using the new
/// interfaces
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag="1")]
    pub params: ::core::option::Option<Params>,
    #[prost(uint64, tag="2")]
    pub last_observed_nonce: u64,
    #[prost(message, repeated, tag="3")]
    pub valsets: ::prost::alloc::vec::Vec<UpdateSignerSetTx>,
    #[prost(message, repeated, tag="4")]
    pub valset_confirms: ::prost::alloc::vec::Vec<UpdateSignerSetTxSignature>,
    #[prost(message, repeated, tag="5")]
    pub batches: ::prost::alloc::vec::Vec<BatchTx>,
    #[prost(message, repeated, tag="6")]
    pub batch_confirms: ::prost::alloc::vec::Vec<BatchTxSignature>,
    #[prost(message, repeated, tag="7")]
    pub logic_calls: ::prost::alloc::vec::Vec<ContractCallTx>,
    #[prost(message, repeated, tag="8")]
    pub logic_call_confirms: ::prost::alloc::vec::Vec<ContractCallTxSignature>,
    #[prost(message, repeated, tag="9")]
    pub ethereum_event_vote_records: ::prost::alloc::vec::Vec<EthereumEventVoteRecord>,
    #[prost(message, repeated, tag="10")]
    pub delegate_keys: ::prost::alloc::vec::Vec<MsgDelegateKeys>,
    #[prost(message, repeated, tag="11")]
    pub erc20_to_denoms: ::prost::alloc::vec::Vec<Erc20ToDenom>,
    #[prost(message, repeated, tag="12")]
    pub unbatched_send_to_ethereum_txs: ::prost::alloc::vec::Vec<SendToEthereumTx>,
}
/// This records the relationship between an ERC20 token and the denom
/// of the corresponding Cosmos originated asset
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20ToDenom {
    #[prost(string, tag="1")]
    pub erc20: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub denom: ::prost::alloc::string::String,
}
/// rpc Params
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ParamsRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ParamsResponse {
    #[prost(message, optional, tag="1")]
    pub params: ::core::option::Option<Params>,
}
/// rpc UpdateSignerSetTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateSignerSetTxRequest {
    /// NOTE: if nonce is not passed, then return the current
    #[prost(uint64, tag="1")]
    pub nonce: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateSignerSetTxResponse {
    #[prost(message, optional, tag="1")]
    pub signer_set: ::core::option::Option<UpdateSignerSetTx>,
}
/// rpc BatchTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxRequest {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(string, tag="2")]
    pub contract_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxResponse {
    #[prost(message, optional, tag="1")]
    pub batch: ::core::option::Option<BatchTx>,
}
/// rpc ContractCallTx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxRequest {
    #[prost(bytes="vec", tag="1")]
    pub invalidation_id: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="2")]
    pub invalidation_nonce: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxResponse {
    #[prost(message, optional, tag="1")]
    pub logic_call: ::core::option::Option<ContractCallTx>,
}
/// rpc UpdateSignerSetTxEthereumSignatures
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateSignerSetTxEthereumSignaturesRequest {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    /// NOTE: if address is passed, return only the signature from that validator
    #[prost(string, tag="2")]
    pub address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateSignerSetTxEthereumSignaturesResponse {
    #[prost(message, repeated, tag="1")]
    pub confirm: ::prost::alloc::vec::Vec<UpdateSignerSetTxSignature>,
}
/// rpc UpdateSignerSetTxs
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateSignerSetTxsRequest {
    #[prost(int64, tag="1")]
    pub count: i64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateSignerSetTxsResponse {
    #[prost(message, repeated, tag="1")]
    pub signer_sets: ::prost::alloc::vec::Vec<UpdateSignerSetTx>,
}
/// rpc PendingUpdateSignerSetTxEthereumSignatures
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PendingUpdateSignerSetTxEthereumSignaturesRequest {
    /// NOTE: this is an sdk.AccAddress and can represent either the 
    /// orchestartor address or the cooresponding validator address
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PendingUpdateSignerSetTxEthereumSignaturesResponse {
    #[prost(message, repeated, tag="1")]
    pub signer_sets: ::prost::alloc::vec::Vec<UpdateSignerSetTx>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PendingBatchTxEthereumSignaturesRequest {
    /// NOTE: this is an sdk.AccAddress and can represent either the 
    /// orchestartor address or the cooresponding validator address
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PendingBatchTxEthereumSignaturesResponse {
    /// Note these are returned with the signature empty
    #[prost(message, repeated, tag="1")]
    pub batches: ::prost::alloc::vec::Vec<BatchTx>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PendingContractCallTxEthereumSignaturesRequest {
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PendingContractCallTxEthereumSignaturesResponse {
    #[prost(message, repeated, tag="1")]
    pub calls: ::prost::alloc::vec::Vec<ContractCallTx>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxFeesRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxFeesResponse {
    #[prost(message, repeated, tag="1")]
    pub fees: ::prost::alloc::vec::Vec<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxEthereumSignaturesRequest {
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxEthereumSignaturesResponse {
    /// Note these are returned with the signature empty
    #[prost(message, repeated, tag="1")]
    pub logic_call_confirms: ::prost::alloc::vec::Vec<ContractCallTx>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxsRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxsResponse {
    #[prost(message, repeated, tag="1")]
    pub batches: ::prost::alloc::vec::Vec<BatchTx>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxsRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ContractCallTxsResponse {
    #[prost(message, repeated, tag="1")]
    pub calls: ::prost::alloc::vec::Vec<ContractCallTx>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxEthereumSignaturesRequest {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(string, tag="2")]
    pub contract_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchTxEthereumSignaturesResponse {
    #[prost(message, repeated, tag="1")]
    pub confirms: ::prost::alloc::vec::Vec<BatchTxSignature>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LastSubmittedEthereumEventRequest {
    #[prost(string, tag="1")]
    pub address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LastSubmittedEthereumEventResponse {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20ToDenomRequest {
    #[prost(string, tag="1")]
    pub erc20: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20ToDenomResponse {
    #[prost(string, tag="1")]
    pub denom: ::prost::alloc::string::String,
    #[prost(bool, tag="2")]
    pub cosmos_originated: bool,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DenomToErc20Request {
    #[prost(string, tag="1")]
    pub denom: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DenomToErc20Response {
    #[prost(string, tag="1")]
    pub erc20: ::prost::alloc::string::String,
    #[prost(bool, tag="2")]
    pub cosmos_originated: bool,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByValidatorAddress {
    #[prost(string, tag="1")]
    pub validator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByValidatorAddressResponse {
    #[prost(string, tag="1")]
    pub eth_address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByEthereumSignerRequest {
    #[prost(string, tag="1")]
    pub ethereum_signer: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByEthereumSignerResponse {
    #[prost(string, tag="1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByOrchestratorAddress {
    #[prost(string, tag="1")]
    pub orchestrator_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DelegateKeysByOrchestratorAddressResponse {
    #[prost(string, tag="1")]
    pub validator_address: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub ethereum_signer: ::prost::alloc::string::String,
}
/// NOTE: if there is no sender address, return all
/// TODO: pagination
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PendingSendToEthereumTxsRequest {
    #[prost(string, tag="1")]
    pub sender_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PendingSendToEthereumTxsResponse {
    #[prost(message, repeated, tag="1")]
    pub batched_send_to_ethereum_txs: ::prost::alloc::vec::Vec<SendToEthereumTx>,
    #[prost(message, repeated, tag="2")]
    pub unbatched_send_to_ethereum_txs: ::prost::alloc::vec::Vec<SendToEthereumTx>,
}
# [doc = r" Generated client implementations."] pub mod query_client { # ! [allow (unused_variables , dead_code , missing_docs)] use tonic :: codegen :: * ; # [doc = " Query defines the gRPC querier service"] pub struct QueryClient < T > { inner : tonic :: client :: Grpc < T > , } impl QueryClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > QueryClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + HttpBody + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as HttpBody > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor (inner : T , interceptor : impl Into < tonic :: Interceptor >) -> Self { let inner = tonic :: client :: Grpc :: with_interceptor (inner , interceptor) ; Self { inner } } # [doc = " Module parameters query"] pub async fn params (& mut self , request : impl tonic :: IntoRequest < super :: ParamsRequest > ,) -> Result < tonic :: Response < super :: ParamsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/Params") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " get info on individual outgoing data"] pub async fn update_signer_set_tx (& mut self , request : impl tonic :: IntoRequest < super :: UpdateSignerSetTxRequest > ,) -> Result < tonic :: Response < super :: UpdateSignerSetTxResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/UpdateSignerSetTx") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn batch_tx (& mut self , request : impl tonic :: IntoRequest < super :: BatchTxRequest > ,) -> Result < tonic :: Response < super :: BatchTxResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/BatchTx") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn contract_call_tx (& mut self , request : impl tonic :: IntoRequest < super :: ContractCallTxRequest > ,) -> Result < tonic :: Response < super :: ContractCallTxResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/ContractCallTx") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " get collections of outgoing traffic from the bridge"] pub async fn update_signer_set_txs (& mut self , request : impl tonic :: IntoRequest < super :: UpdateSignerSetTxsRequest > ,) -> Result < tonic :: Response < super :: UpdateSignerSetTxsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/UpdateSignerSetTxs") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn batch_txs (& mut self , request : impl tonic :: IntoRequest < super :: BatchTxsRequest > ,) -> Result < tonic :: Response < super :: BatchTxsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/BatchTxs") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn contract_call_txs (& mut self , request : impl tonic :: IntoRequest < super :: ContractCallTxsRequest > ,) -> Result < tonic :: Response < super :: ContractCallTxsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/ContractCallTxs") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " ethereum signature queries so validators can construct valid etherum transactions"] # [doc = " TODO: can/should we group these into one endpoint?"] pub async fn update_signer_set_tx_ethereum_signatures (& mut self , request : impl tonic :: IntoRequest < super :: UpdateSignerSetTxEthereumSignaturesRequest > ,) -> Result < tonic :: Response < super :: UpdateSignerSetTxEthereumSignaturesResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/UpdateSignerSetTxEthereumSignatures") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn batch_tx_ethereum_signatures (& mut self , request : impl tonic :: IntoRequest < super :: BatchTxEthereumSignaturesRequest > ,) -> Result < tonic :: Response < super :: BatchTxEthereumSignaturesResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/BatchTxEthereumSignatures") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn contract_call_tx_ethereum_signatures (& mut self , request : impl tonic :: IntoRequest < super :: ContractCallTxEthereumSignaturesRequest > ,) -> Result < tonic :: Response < super :: ContractCallTxEthereumSignaturesResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/ContractCallTxEthereumSignatures") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " pending ethereum signature queries for orchestrators to figure out which signatures they are missing"] # [doc = " TODO: can/should we group this into one endpoint?"] pub async fn pending_update_signer_set_tx_ethereum_signatures (& mut self , request : impl tonic :: IntoRequest < super :: PendingUpdateSignerSetTxEthereumSignaturesRequest > ,) -> Result < tonic :: Response < super :: PendingUpdateSignerSetTxEthereumSignaturesResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/PendingUpdateSignerSetTxEthereumSignatures") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn pending_batch_tx_ethereum_signatures (& mut self , request : impl tonic :: IntoRequest < super :: PendingBatchTxEthereumSignaturesRequest > ,) -> Result < tonic :: Response < super :: PendingBatchTxEthereumSignaturesResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/PendingBatchTxEthereumSignatures") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn pending_contract_call_tx_ethereum_signatures (& mut self , request : impl tonic :: IntoRequest < super :: PendingContractCallTxEthereumSignaturesRequest > ,) -> Result < tonic :: Response < super :: PendingContractCallTxEthereumSignaturesResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/PendingContractCallTxEthereumSignatures") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn last_submitted_ethereum_event (& mut self , request : impl tonic :: IntoRequest < super :: LastSubmittedEthereumEventRequest > ,) -> Result < tonic :: Response < super :: LastSubmittedEthereumEventResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/LastSubmittedEthereumEvent") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Queries the fees for all pending batches, results are returned in sdk.Coin (fee_amount_int)(contract_address) style"] pub async fn batch_tx_fees (& mut self , request : impl tonic :: IntoRequest < super :: BatchTxFeesRequest > ,) -> Result < tonic :: Response < super :: BatchTxFeesResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/BatchTxFees") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Query for info about denoms tracked by gravity"] pub async fn erc20_to_denom (& mut self , request : impl tonic :: IntoRequest < super :: Erc20ToDenomRequest > ,) -> Result < tonic :: Response < super :: Erc20ToDenomResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/ERC20ToDenom") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Query for info about denoms tracked by gravity"] pub async fn denom_to_erc20 (& mut self , request : impl tonic :: IntoRequest < super :: DenomToErc20Request > ,) -> Result < tonic :: Response < super :: DenomToErc20Response > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/DenomToERC20") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " Query for pending tranfertxs"] pub async fn pending_send_to_ethereum_txs (& mut self , request : impl tonic :: IntoRequest < super :: PendingSendToEthereumTxsRequest > ,) -> Result < tonic :: Response < super :: PendingSendToEthereumTxsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/PendingSendToEthereumTxs") ; self . inner . unary (request . into_request () , path , codec) . await } # [doc = " delegate keys"] pub async fn delegate_keys_by_validator (& mut self , request : impl tonic :: IntoRequest < super :: DelegateKeysByValidatorAddress > ,) -> Result < tonic :: Response < super :: DelegateKeysByValidatorAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/DelegateKeysByValidator") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn delegate_keys_by_ethereum_signer (& mut self , request : impl tonic :: IntoRequest < super :: DelegateKeysByEthereumSignerRequest > ,) -> Result < tonic :: Response < super :: DelegateKeysByEthereumSignerResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/DelegateKeysByEthereumSigner") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn delegate_keys_by_orchestrator (& mut self , request : impl tonic :: IntoRequest < super :: DelegateKeysByOrchestratorAddress > ,) -> Result < tonic :: Response < super :: DelegateKeysByOrchestratorAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/gravity.v1.Query/DelegateKeysByOrchestrator") ; self . inner . unary (request . into_request () , path , codec) . await } } impl < T : Clone > Clone for QueryClient < T > { fn clone (& self) -> Self { Self { inner : self . inner . clone () , } } } impl < T > std :: fmt :: Debug for QueryClient < T > { fn fmt (& self , f : & mut std :: fmt :: Formatter < '_ >) -> std :: fmt :: Result { write ! (f , "QueryClient {{ ... }}") } } }
