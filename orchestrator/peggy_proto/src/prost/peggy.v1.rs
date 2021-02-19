/// Attestation is an aggregate of `claims` that eventually becomes `observed` by
/// all orchestrators
/// EVENT_NONCE:
/// EventNonce a nonce provided by the peggy contract that is unique per event fired
/// These event nonces must be relayed in order. This is a correctness issue,
/// if relaying out of order transaction replay attacks become possible
/// OBSERVED:
/// Observed indicates that >67% of validators have attested to the event,
/// and that the event should be executed by the peggy state machine
///
/// The actual content of the claims is passed in with the transaction making the claim
/// and then passed through the call stack alongside the attestation while it is processed
/// the key in which the attestation is stored is keyed on the exact details of the claim
/// but there is no reason to store those exact details becuause the next message sender
/// will kindly provide you with them.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Attestation {
    #[prost(bool, tag="1")]
    pub observed: bool,
    #[prost(string, repeated, tag="2")]
    pub votes: ::std::vec::Vec<std::string::String>,
    #[prost(uint64, tag="3")]
    pub height: u64,
    #[prost(message, optional, tag="4")]
    pub claim: ::std::option::Option<::prost_types::Any>,
}
/// ERC20Token unique identifier for an Ethereum ERC20 token.
/// CONTRACT:
/// The contract address on ETH of the token (note: developers should look up
/// the token symbol using the address on ETH to display for UI)
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Erc20Token {
    #[prost(string, tag="1")]
    pub contract: std::string::String,
    #[prost(string, tag="2")]
    pub amount: std::string::String,
}
/// ClaimType is the cosmos type of an event from the counterpart chain that can
/// be handled
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum ClaimType {
    Unknown = 0,
    Deposit = 1,
    Withdraw = 2,
    Erc20Deployed = 3,
    LogicCallExecuted = 4,
}
/// OutgoingTx is a withdrawal on the bridged contract
/// TODO: can this type be replaced by outgoing transfer tx
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct OutgoingTx {
    #[prost(string, tag="1")]
    pub sender: std::string::String,
    #[prost(string, tag="2")]
    pub dest_addr: std::string::String,
    #[prost(message, optional, tag="3")]
    pub amount: ::std::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, optional, tag="4")]
    pub bridge_fee: ::std::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
/// IDSet represents a set of IDs
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct IdSet {
    #[prost(uint64, repeated, tag="1")]
    pub ids: ::std::vec::Vec<u64>,
}
/// OutgoingTxBatch represents a batch of transactions going from Peggy to ETH
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct OutgoingTxBatch {
    #[prost(uint64, tag="1")]
    pub batch_nonce: u64,
    #[prost(uint64, tag="2")]
    pub batch_timeout: u64,
    #[prost(message, repeated, tag="3")]
    pub transactions: ::std::vec::Vec<OutgoingTransferTx>,
    #[prost(string, tag="4")]
    pub token_contract: std::string::String,
    #[prost(uint64, tag="5")]
    pub block: u64,
}
/// OutgoingTransferTx represents an individual send from Peggy to ETH
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct OutgoingTransferTx {
    #[prost(uint64, tag="1")]
    pub id: u64,
    #[prost(string, tag="2")]
    pub sender: std::string::String,
    #[prost(string, tag="3")]
    pub dest_address: std::string::String,
    #[prost(message, optional, tag="4")]
    pub erc20_token: ::std::option::Option<Erc20Token>,
    #[prost(message, optional, tag="5")]
    pub erc20_fee: ::std::option::Option<Erc20Token>,
}
/// OutgoingLogicCall represents an individual logic call from Peggy to ETH
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct OutgoingLogicCall {
    #[prost(message, repeated, tag="1")]
    pub transfers: ::std::vec::Vec<Erc20Token>,
    #[prost(message, repeated, tag="2")]
    pub fees: ::std::vec::Vec<Erc20Token>,
    #[prost(string, tag="3")]
    pub logic_contract_address: std::string::String,
    #[prost(bytes, tag="4")]
    pub payload: std::vec::Vec<u8>,
    #[prost(uint64, tag="5")]
    pub timeout: u64,
    #[prost(bytes, tag="6")]
    pub invalidation_id: std::vec::Vec<u8>,
    #[prost(uint64, tag="7")]
    pub invalidation_nonce: u64,
}
/// SignType defines messages that have been signed by an orchestrator
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum SignType {
    Unknown = 0,
    OrchestratorSignedMultiSigUpdate = 1,
    OrchestratorSignedWithdrawBatch = 2,
}
/// BridgeValidator represents a validator's ETH address and its power
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BridgeValidator {
    #[prost(uint64, tag="1")]
    pub power: u64,
    #[prost(string, tag="2")]
    pub ethereum_address: std::string::String,
}
/// Valset is the Ethereum Bridge Multsig Set, each peggy validator also
/// maintains an ETH key to sign messages, these are used to check signatures on
/// ETH because of the significant gas savings
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Valset {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(message, repeated, tag="2")]
    pub members: ::std::vec::Vec<BridgeValidator>,
    #[prost(uint64, tag="3")]
    pub height: u64,
}
/// LastObservedEthereumBlockHeight stores the last observed
/// Ethereum block height along with the Cosmos block height that
/// it was observed at. These two numbers can be used to project
/// outward and always produce batches with timeouts in the future
/// even if no Ethereum block height has been relayed for a long time
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LastObservedEthereumBlockHeight {
    #[prost(uint64, tag="1")]
    pub cosmos_block_height: u64,
    #[prost(uint64, tag="2")]
    pub ethereum_block_height: u64,
}
/// MsgSetOrchestratorAddress
/// this message allows validators to delegate their voting responsibilities 
/// to a given key. This key is then used as an optional authentication method
/// for sigining oracle claims
/// VALIDATOR
/// The validator field is a cosmosvaloper1... string (i.e. sdk.ValAddress)
/// that references a validator in the active set
/// ORCHESTRATOR
/// The orchestrator field is a cosmos1... string  (i.e. sdk.AccAddress) that 
/// references the key that is being delegated to
/// ETH_ADDRESS
/// This is a hex encoded 0x Ethereum public key that will be used by this validator
/// on Ethereum
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSetOrchestratorAddress {
    #[prost(string, tag="1")]
    pub validator: std::string::String,
    #[prost(string, tag="2")]
    pub orchestrator: std::string::String,
    #[prost(string, tag="3")]
    pub eth_address: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSetOrchestratorAddressResponse {
}
/// MsgValsetConfirm
/// this is the message sent by the validators when they wish to submit their
/// signatures over the validator set at a given block height. A validator must
/// first call MsgSetEthAddress to set their Ethereum address to be used for
/// signing. Then someone (anyone) must make a ValsetRequest the request is
/// essentially a messaging mechanism to determine which block all validators
/// should submit signatures over. Finally validators sign the validator set,
/// powers, and Ethereum addresses of the entire validator set at the height of a
/// ValsetRequest and submit that signature with this message.
///
/// If a sufficient number of validators (66% of voting power) (A) have set
/// Ethereum addresses and (B) submit ValsetConfirm messages with their
/// signatures it is then possible for anyone to view these signatures in the
/// chain store and submit them to Ethereum to update the validator set
/// -------------
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgValsetConfirm {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(string, tag="2")]
    pub orchestrator: std::string::String,
    #[prost(string, tag="3")]
    pub eth_address: std::string::String,
    #[prost(string, tag="4")]
    pub signature: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgValsetConfirmResponse {
}
/// MsgSendToEth
/// This is the message that a user calls when they want to bridge an asset
/// it will later be removed when it is included in a batch and successfully
/// submitted tokens are removed from the users balance immediately
/// -------------
/// AMOUNT:
/// the coin to send across the bridge, note the restriction that this is a
/// single coin not a set of coins that is normal in other Cosmos messages
/// FEE:
/// the fee paid for the bridge, distinct from the fee paid to the chain to
/// actually send this message in the first place. So a successful send has
/// two layers of fees for the user
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSendToEth {
    #[prost(string, tag="1")]
    pub sender: std::string::String,
    #[prost(string, tag="2")]
    pub eth_dest: std::string::String,
    #[prost(message, optional, tag="3")]
    pub amount: ::std::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, optional, tag="4")]
    pub bridge_fee: ::std::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSendToEthResponse {
}
/// MsgRequestBatch
/// this is a message anyone can send that requests a batch of transactions to
/// send across the bridge be created for whatever block height this message is
/// included in. This acts as a coordination point, the handler for this message
/// looks at the AddToOutgoingPool tx's in the store and generates a batch, also
/// available in the store tied to this message. The validators then grab this
/// batch, sign it, submit the signatures with a MsgConfirmBatch before a relayer
/// can finally submit the batch
/// -------------
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRequestBatch {
    #[prost(string, tag="1")]
    pub orchestrator: std::string::String,
    #[prost(string, tag="2")]
    pub denom: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgRequestBatchResponse {
}
/// MsgConfirmBatch
/// When validators observe a MsgRequestBatch they form a batch by ordering
/// transactions currently in the txqueue in order of highest to lowest fee,
/// cutting off when the batch either reaches a hardcoded maximum size (to be
/// decided, probably around 100) or when transactions stop being profitable
/// (TODO determine this without nondeterminism) This message includes the batch
/// as well as an Ethereum signature over this batch by the validator
/// -------------
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgConfirmBatch {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(string, tag="2")]
    pub token_contract: std::string::String,
    #[prost(string, tag="3")]
    pub eth_signer: std::string::String,
    #[prost(string, tag="4")]
    pub orchestrator: std::string::String,
    #[prost(string, tag="5")]
    pub signature: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgConfirmBatchResponse {
}
/// MsgConfirmLogicCall
/// When validators observe a MsgRequestBatch they form a batch by ordering
/// transactions currently in the txqueue in order of highest to lowest fee,
/// cutting off when the batch either reaches a hardcoded maximum size (to be
/// decided, probably around 100) or when transactions stop being profitable
/// (TODO determine this without nondeterminism) This message includes the batch
/// as well as an Ethereum signature over this batch by the validator
/// -------------
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgConfirmLogicCall {
    #[prost(bytes, tag="1")]
    pub invalidation_id: std::vec::Vec<u8>,
    #[prost(uint64, tag="2")]
    pub invalidation_nonce: u64,
    #[prost(string, tag="3")]
    pub eth_signer: std::string::String,
    #[prost(string, tag="4")]
    pub orchestrator: std::string::String,
    #[prost(string, tag="5")]
    pub signature: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgConfirmLogicCallResponse {
}
/// EthereumBridgeDepositClaim
/// When more than 66% of the active validator set has
/// claimed to have seen the deposit enter the ethereum blockchain coins are
/// issued to the Cosmos address in question
/// -------------
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgDepositClaim {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
    #[prost(uint64, tag="2")]
    pub block_height: u64,
    #[prost(string, tag="3")]
    pub token_contract: std::string::String,
    #[prost(string, tag="4")]
    pub amount: std::string::String,
    #[prost(string, tag="5")]
    pub ethereum_sender: std::string::String,
    #[prost(string, tag="6")]
    pub cosmos_receiver: std::string::String,
    #[prost(string, tag="7")]
    pub orchestrator: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgDepositClaimResponse {
}
/// WithdrawClaim claims that a batch of withdrawal
/// operations on the bridge contract was executed.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgWithdrawClaim {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
    #[prost(uint64, tag="2")]
    pub block_height: u64,
    #[prost(uint64, tag="3")]
    pub batch_nonce: u64,
    #[prost(string, tag="4")]
    pub token_contract: std::string::String,
    #[prost(string, tag="5")]
    pub orchestrator: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgWithdrawClaimResponse {
}
/// ERC20DeployedClaim allows the Cosmos module
/// to learn about an ERC20 that someone deployed
/// to represent a Cosmos asset
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgErc20DeployedClaim {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
    #[prost(uint64, tag="2")]
    pub block_height: u64,
    #[prost(string, tag="3")]
    pub cosmos_denom: std::string::String,
    #[prost(string, tag="4")]
    pub token_contract: std::string::String,
    #[prost(string, tag="5")]
    pub name: std::string::String,
    #[prost(string, tag="6")]
    pub symbol: std::string::String,
    #[prost(uint64, tag="7")]
    pub decimals: u64,
    #[prost(string, tag="8")]
    pub orchestrator: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgErc20DeployedClaimResponse {
}
/// This informs the Cosmos module that a logic
/// call has been executed
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgLogicCallExecutedClaim {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
    #[prost(uint64, tag="2")]
    pub block_height: u64,
    #[prost(bytes, tag="3")]
    pub invalidation_id: std::vec::Vec<u8>,
    #[prost(uint64, tag="4")]
    pub invalidation_nonce: u64,
    #[prost(string, tag="5")]
    pub orchestrator: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgLogicCallExecutedClaimResponse {
}
# [doc = r" Generated client implementations."] pub mod msg_client { # ! [allow (unused_variables , dead_code , missing_docs)] use tonic :: codegen :: * ; pub struct MsgClient < T > { inner : tonic :: client :: Grpc < T > , } impl MsgClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > MsgClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + HttpBody + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as HttpBody > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor (inner : T , interceptor : impl Into < tonic :: Interceptor >) -> Self { let inner = tonic :: client :: Grpc :: with_interceptor (inner , interceptor) ; Self { inner } } pub async fn valset_confirm (& mut self , request : impl tonic :: IntoRequest < super :: MsgValsetConfirm > ,) -> Result < tonic :: Response < super :: MsgValsetConfirmResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/ValsetConfirm") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn send_to_eth (& mut self , request : impl tonic :: IntoRequest < super :: MsgSendToEth > ,) -> Result < tonic :: Response < super :: MsgSendToEthResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/SendToEth") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn request_batch (& mut self , request : impl tonic :: IntoRequest < super :: MsgRequestBatch > ,) -> Result < tonic :: Response < super :: MsgRequestBatchResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/RequestBatch") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn confirm_batch (& mut self , request : impl tonic :: IntoRequest < super :: MsgConfirmBatch > ,) -> Result < tonic :: Response < super :: MsgConfirmBatchResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/ConfirmBatch") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn confirm_logic_call (& mut self , request : impl tonic :: IntoRequest < super :: MsgConfirmLogicCall > ,) -> Result < tonic :: Response < super :: MsgConfirmLogicCallResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/ConfirmLogicCall") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn deposit_claim (& mut self , request : impl tonic :: IntoRequest < super :: MsgDepositClaim > ,) -> Result < tonic :: Response < super :: MsgDepositClaimResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/DepositClaim") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn withdraw_claim (& mut self , request : impl tonic :: IntoRequest < super :: MsgWithdrawClaim > ,) -> Result < tonic :: Response < super :: MsgWithdrawClaimResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/WithdrawClaim") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn erc20_deployed_claim (& mut self , request : impl tonic :: IntoRequest < super :: MsgErc20DeployedClaim > ,) -> Result < tonic :: Response < super :: MsgErc20DeployedClaimResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/ERC20DeployedClaim") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn logic_call_executed_claim (& mut self , request : impl tonic :: IntoRequest < super :: MsgLogicCallExecutedClaim > ,) -> Result < tonic :: Response < super :: MsgLogicCallExecutedClaimResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/LogicCallExecutedClaim") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn set_orchestrator_address (& mut self , request : impl tonic :: IntoRequest < super :: MsgSetOrchestratorAddress > ,) -> Result < tonic :: Response < super :: MsgSetOrchestratorAddressResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Msg/SetOrchestratorAddress") ; self . inner . unary (request . into_request () , path , codec) . await } } impl < T : Clone > Clone for MsgClient < T > { fn clone (& self) -> Self { Self { inner : self . inner . clone () , } } } impl < T > std :: fmt :: Debug for MsgClient < T > { fn fmt (& self , f : & mut std :: fmt :: Formatter < '_ >) -> std :: fmt :: Result { write ! (f , "MsgClient {{ ... }}") } } }/// Params represent the peggy genesis and store parameters
/// peggy_id:
/// a random 32 byte value to prevent signature reuse, for example if the
/// cosmos validators decided to use the same Ethereum keys for another chain
/// also running Peggy we would not want it to be possible to play a deposit
/// from chain A back on chain B's peggy. This value IS USED ON ETHEREUM so
/// it must be set in your genesis.json before launch and not changed after
/// deploying Peggy
///
/// contract_hash:
/// the code hash of a known good version of the Peggy contract
/// solidity code. This can be used to verify the correct version
/// of the contract has been deployed. This is a reference value for
/// goernance action only it is never read by any Peggy code
///
/// bridge_ethereum_address:
/// is address of the bridge contract on the Ethereum side, this is a
/// reference value for governance only and is not actually used by any
/// Peggy code
///
/// bridge_chain_id:
/// the unique identifier of the Ethereum chain, this is a reference value
/// only and is not actually used by any Peggy code
///
/// These reference values may be used by future Peggy client implemetnations
/// to allow for saftey features or convenience features like the peggy address
/// in your relayer. A relayer would require a configured peggy address if
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
/// The slashing fractions for the various peggy related slashing conditions. The first three
/// refer to not submitting a particular message, the third for submitting a different claim
/// for the same Ethereum event
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
    #[prost(string, tag="1")]
    pub peggy_id: std::string::String,
    #[prost(string, tag="2")]
    pub contract_source_hash: std::string::String,
    #[prost(string, tag="4")]
    pub bridge_ethereum_address: std::string::String,
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
    #[prost(bytes, tag="13")]
    pub slash_fraction_valset: std::vec::Vec<u8>,
    #[prost(bytes, tag="14")]
    pub slash_fraction_batch: std::vec::Vec<u8>,
    #[prost(bytes, tag="15")]
    pub slash_fraction_claim: std::vec::Vec<u8>,
    #[prost(bytes, tag="16")]
    pub slash_fraction_conflicting_claim: std::vec::Vec<u8>,
}
/// GenesisState struct
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag="1")]
    pub params: ::std::option::Option<Params>,
    #[prost(uint64, tag="2")]
    pub last_observed_nonce: u64,
    #[prost(message, repeated, tag="3")]
    pub valsets: ::std::vec::Vec<Valset>,
    #[prost(message, repeated, tag="4")]
    pub valset_confirms: ::std::vec::Vec<MsgValsetConfirm>,
    #[prost(message, repeated, tag="5")]
    pub batches: ::std::vec::Vec<OutgoingTxBatch>,
    #[prost(message, repeated, tag="6")]
    pub batch_confirms: ::std::vec::Vec<MsgConfirmBatch>,
    #[prost(message, repeated, tag="7")]
    pub logic_calls: ::std::vec::Vec<OutgoingLogicCall>,
    #[prost(message, repeated, tag="8")]
    pub logic_call_confirms: ::std::vec::Vec<MsgConfirmLogicCall>,
    #[prost(message, repeated, tag="9")]
    pub attestations: ::std::vec::Vec<Attestation>,
    #[prost(message, repeated, tag="10")]
    pub delegate_keys: ::std::vec::Vec<MsgSetOrchestratorAddress>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    #[prost(message, optional, tag="1")]
    pub params: ::std::option::Option<Params>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCurrentValsetRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCurrentValsetResponse {
    #[prost(message, optional, tag="1")]
    pub valset: ::std::option::Option<Valset>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryValsetRequestRequest {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryValsetRequestResponse {
    #[prost(message, optional, tag="1")]
    pub valset: ::std::option::Option<Valset>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryValsetConfirmRequest {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(string, tag="2")]
    pub address: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryValsetConfirmResponse {
    #[prost(message, optional, tag="1")]
    pub confirm: ::std::option::Option<MsgValsetConfirm>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryValsetConfirmsByNonceRequest {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryValsetConfirmsByNonceResponse {
    #[prost(message, repeated, tag="1")]
    pub confirms: ::std::vec::Vec<MsgValsetConfirm>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastValsetRequestsRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastValsetRequestsResponse {
    #[prost(message, repeated, tag="1")]
    pub valsets: ::std::vec::Vec<Valset>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastPendingValsetRequestByAddrRequest {
    #[prost(string, tag="1")]
    pub address: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastPendingValsetRequestByAddrResponse {
    #[prost(message, repeated, tag="1")]
    pub valsets: ::std::vec::Vec<Valset>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastPendingBatchRequestByAddrRequest {
    #[prost(string, tag="1")]
    pub address: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastPendingBatchRequestByAddrResponse {
    #[prost(message, optional, tag="1")]
    pub batch: ::std::option::Option<OutgoingTxBatch>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastPendingLogicCallByAddrRequest {
    #[prost(string, tag="1")]
    pub address: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastPendingLogicCallByAddrResponse {
    #[prost(message, optional, tag="1")]
    pub call: ::std::option::Option<OutgoingLogicCall>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryOutgoingTxBatchesRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryOutgoingTxBatchesResponse {
    #[prost(message, repeated, tag="1")]
    pub batches: ::std::vec::Vec<OutgoingTxBatch>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryOutgoingLogicCallsRequest {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryOutgoingLogicCallsResponse {
    #[prost(message, repeated, tag="1")]
    pub calls: ::std::vec::Vec<OutgoingLogicCall>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBatchRequestByNonceRequest {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(string, tag="2")]
    pub contract_address: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBatchRequestByNonceResponse {
    #[prost(message, optional, tag="1")]
    pub batch: ::std::option::Option<OutgoingTxBatch>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBatchConfirmsRequest {
    #[prost(uint64, tag="1")]
    pub nonce: u64,
    #[prost(string, tag="2")]
    pub contract_address: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBatchConfirmsResponse {
    #[prost(message, repeated, tag="1")]
    pub confirms: ::std::vec::Vec<MsgConfirmBatch>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLogicConfirmsRequest {
    #[prost(bytes, tag="1")]
    pub invalidation_id: std::vec::Vec<u8>,
    #[prost(uint64, tag="2")]
    pub invalidation_nonce: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLogicConfirmsResponse {
    #[prost(message, repeated, tag="1")]
    pub confirms: ::std::vec::Vec<MsgConfirmLogicCall>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastEventNonceByAddrRequest {
    #[prost(string, tag="1")]
    pub address: std::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryLastEventNonceByAddrResponse {
    #[prost(uint64, tag="1")]
    pub event_nonce: u64,
}
# [doc = r" Generated client implementations."] pub mod query_client { # ! [allow (unused_variables , dead_code , missing_docs)] use tonic :: codegen :: * ; # [doc = " Query defines the gRPC querier service"] pub struct QueryClient < T > { inner : tonic :: client :: Grpc < T > , } impl QueryClient < tonic :: transport :: Channel > { # [doc = r" Attempt to create a new client by connecting to a given endpoint."] pub async fn connect < D > (dst : D) -> Result < Self , tonic :: transport :: Error > where D : std :: convert :: TryInto < tonic :: transport :: Endpoint > , D :: Error : Into < StdError > , { let conn = tonic :: transport :: Endpoint :: new (dst) ? . connect () . await ? ; Ok (Self :: new (conn)) } } impl < T > QueryClient < T > where T : tonic :: client :: GrpcService < tonic :: body :: BoxBody > , T :: ResponseBody : Body + HttpBody + Send + 'static , T :: Error : Into < StdError > , < T :: ResponseBody as HttpBody > :: Error : Into < StdError > + Send , { pub fn new (inner : T) -> Self { let inner = tonic :: client :: Grpc :: new (inner) ; Self { inner } } pub fn with_interceptor (inner : T , interceptor : impl Into < tonic :: Interceptor >) -> Self { let inner = tonic :: client :: Grpc :: with_interceptor (inner , interceptor) ; Self { inner } } # [doc = " Deployments queries deployments"] pub async fn params (& mut self , request : impl tonic :: IntoRequest < super :: QueryParamsRequest > ,) -> Result < tonic :: Response < super :: QueryParamsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/Params") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn current_valset (& mut self , request : impl tonic :: IntoRequest < super :: QueryCurrentValsetRequest > ,) -> Result < tonic :: Response < super :: QueryCurrentValsetResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/CurrentValset") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn valset_request (& mut self , request : impl tonic :: IntoRequest < super :: QueryValsetRequestRequest > ,) -> Result < tonic :: Response < super :: QueryValsetRequestResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/ValsetRequest") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn valset_confirm (& mut self , request : impl tonic :: IntoRequest < super :: QueryValsetConfirmRequest > ,) -> Result < tonic :: Response < super :: QueryValsetConfirmResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/ValsetConfirm") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn valset_confirms_by_nonce (& mut self , request : impl tonic :: IntoRequest < super :: QueryValsetConfirmsByNonceRequest > ,) -> Result < tonic :: Response < super :: QueryValsetConfirmsByNonceResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/ValsetConfirmsByNonce") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn last_valset_requests (& mut self , request : impl tonic :: IntoRequest < super :: QueryLastValsetRequestsRequest > ,) -> Result < tonic :: Response < super :: QueryLastValsetRequestsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/LastValsetRequests") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn last_pending_valset_request_by_addr (& mut self , request : impl tonic :: IntoRequest < super :: QueryLastPendingValsetRequestByAddrRequest > ,) -> Result < tonic :: Response < super :: QueryLastPendingValsetRequestByAddrResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/LastPendingValsetRequestByAddr") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn last_pending_batch_request_by_addr (& mut self , request : impl tonic :: IntoRequest < super :: QueryLastPendingBatchRequestByAddrRequest > ,) -> Result < tonic :: Response < super :: QueryLastPendingBatchRequestByAddrResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/LastPendingBatchRequestByAddr") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn last_pending_logic_call_by_addr (& mut self , request : impl tonic :: IntoRequest < super :: QueryLastPendingLogicCallByAddrRequest > ,) -> Result < tonic :: Response < super :: QueryLastPendingLogicCallByAddrResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/LastPendingLogicCallByAddr") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn last_event_nonce_by_addr (& mut self , request : impl tonic :: IntoRequest < super :: QueryLastEventNonceByAddrRequest > ,) -> Result < tonic :: Response < super :: QueryLastEventNonceByAddrResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/LastEventNonceByAddr") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn outgoing_tx_batches (& mut self , request : impl tonic :: IntoRequest < super :: QueryOutgoingTxBatchesRequest > ,) -> Result < tonic :: Response < super :: QueryOutgoingTxBatchesResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/OutgoingTxBatches") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn outgoing_logic_calls (& mut self , request : impl tonic :: IntoRequest < super :: QueryOutgoingLogicCallsRequest > ,) -> Result < tonic :: Response < super :: QueryOutgoingLogicCallsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/OutgoingLogicCalls") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn batch_request_by_nonce (& mut self , request : impl tonic :: IntoRequest < super :: QueryBatchRequestByNonceRequest > ,) -> Result < tonic :: Response < super :: QueryBatchRequestByNonceResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/BatchRequestByNonce") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn batch_confirms (& mut self , request : impl tonic :: IntoRequest < super :: QueryBatchConfirmsRequest > ,) -> Result < tonic :: Response < super :: QueryBatchConfirmsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/BatchConfirms") ; self . inner . unary (request . into_request () , path , codec) . await } pub async fn logic_confirms (& mut self , request : impl tonic :: IntoRequest < super :: QueryLogicConfirmsRequest > ,) -> Result < tonic :: Response < super :: QueryLogicConfirmsResponse > , tonic :: Status > { self . inner . ready () . await . map_err (| e | { tonic :: Status :: new (tonic :: Code :: Unknown , format ! ("Service was not ready: {}" , e . into ())) }) ? ; let codec = tonic :: codec :: ProstCodec :: default () ; let path = http :: uri :: PathAndQuery :: from_static ("/peggy.v1.Query/LogicConfirms") ; self . inner . unary (request . into_request () , path , codec) . await } } impl < T : Clone > Clone for QueryClient < T > { fn clone (& self) -> Self { Self { inner : self . inner . clone () , } } } impl < T > std :: fmt :: Debug for QueryClient < T > { fn fmt (& self , f : & mut std :: fmt :: Formatter < '_ >) -> std :: fmt :: Result { write ! (f , "QueryClient {{ ... }}") } } }