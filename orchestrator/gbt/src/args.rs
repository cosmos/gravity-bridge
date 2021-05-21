//! Command line argument definitions for Gravity bridge tools
//! See the clap documentation for how exactly this works, note that doc comments are displayed to the user

use clap::AppSettings;
use clap::Clap;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use deep_space::address::Address as CosmosAddress;
use deep_space::PrivateKey as CosmosPrivateKey;

/// Gravity Bridge tools (gbt) provides tools for interacting with the Althea Gravity bridge for Cosmos based blockchains.
#[derive(Clap)]
#[clap(version = "1.0", author = "Justin Kilpatrick <justin@althea.net>")]
#[clap(setting = AppSettings::ColoredHelp)]
pub struct Opts {
    /// Increase the logging verbosity
    #[clap(short, long)]
    pub verbose: bool,
    /// Decrease the logging verbosity
    #[clap(short, long)]
    pub quiet: bool,
    /// Set the address prefix for the Cosmos chain
    #[clap(short, long)]
    pub address_prefix: Option<String>,
    #[clap(subcommand)]
    pub subcmd: SubCommand,
}

#[derive(Clap)]
pub enum SubCommand {
    Orchestrator(OrchestratorOpts),
    Relayer(RelayerOpts),
    Client(ClientOpts),
    Keys(KeyOpts),
}
/// The Gravity Bridge orchestrator is required for all validators of the Cosmos chain running
/// the Gravity Bridge module. It contains an Ethereum Signer, Oracle, and optional relayer
#[derive(Clap)]
#[clap(setting = AppSettings::ColoredHelp)]
pub struct OrchestratorOpts {}

/// The Gravity Bridge Relayer is an unpermissioned role that takes data from the Cosmos blockchain
/// packages it into Ethereum transactions and is paid to submit these transactions to the Ethereum blockchain
#[derive(Clap)]
#[clap(setting = AppSettings::ColoredHelp)]
pub struct RelayerOpts {}

/// The Gravity Bridge client contains helpful command line tools for interacting with the Gravity bridge
#[derive(Clap)]
#[clap(setting = AppSettings::ColoredHelp)]
pub struct ClientOpts {
    #[clap(subcommand)]
    pub subcmd: ClientSubcommand,
}

#[derive(Clap)]
pub enum ClientSubcommand {
    CosmosToEth(CosmosToEthOpts),
    EthToCosmos(EthToCosmosOpts),
    DeployErc20Representation(DeployErc20RepresentationOpts),
}

/// Send Cosmos tokens to Ethereum
#[derive(Clap)]
#[clap(setting = AppSettings::ColoredHelp)]
pub struct CosmosToEthOpts {
    /// Cosmos phrase
    #[clap(short, long)]
    pub cosmos_phrase: String,
    /// (Optional) The Cosmos gRPC server that will be used to submit the transaction
    #[clap(short, long, default_value = "http://localhost:9090")]
    pub cosmos_grpc: String,
    /// The Denom from the Cosmos chain to bridge
    #[clap(short, long)]
    pub cosmos_denom: String,
    /// The amount of tokens you are sending eg. 1.2 ATOM
    #[clap(short, long, parse(try_from_str))]
    pub amount: f64,
    /// The destination address on the Ethereum chain
    #[clap(short, long, parse(try_from_str))]
    pub eth_destination: EthAddress,
    /// If this command should request a batch to push
    /// your tx along immediately
    #[clap(short, long)]
    pub no_batch: bool,
}

/// Send an Ethereum ERC20 token to Cosmos
#[derive(Clap)]
#[clap(setting = AppSettings::ColoredHelp)]
pub struct EthToCosmosOpts {
    /// The Ethereum private key to register, will be generated if not provided
    #[clap(short, long, parse(try_from_str))]
    pub ethereum_key: EthPrivateKey,
    /// (Optional) The Ethereum RPC server that will be used to submit the transaction
    #[clap(short, long, default_value = "http://localhost:8545")]
    pub ethereum_rpc: String,
    /// The address fo the Gravity contract on Ethereum
    #[clap(short, long, parse(try_from_str))]
    pub gravity_contract_address: EthAddress,
    /// The ERC20 contract address of the ERC20 you are sending
    #[clap(short, long, parse(try_from_str))]
    pub erc20_contract_address: EthAddress,
    /// The amount of tokens you are sending eg. 1.2 ATOM
    #[clap(short, long, parse(try_from_str))]
    pub amount: f64,
    /// The destination address on the Cosmos blockchain
    #[clap(short, long, parse(try_from_str))]
    pub destination: CosmosAddress,
}

/// Deploy an ERC20 representation of a Cosmos asset on the Ethereum chain
/// this can only be run once for each time of Cosmos asset
#[derive(Clap)]
#[clap(setting = AppSettings::ColoredHelp)]
pub struct DeployErc20RepresentationOpts {
    /// (Optional) The Cosmos gRPC server that will be used to submit the transaction
    #[clap(short, long, default_value = "http://localhost:9090")]
    pub cosmos_grpc: String,
    /// (Optional) The Ethereum RPC server that will be used to submit the transaction
    #[clap(short, long, default_value = "http://localhost:8545")]
    pub ethereum_rpc: String,
    /// The Cosmos Denom you wish to create an ERC20 representation for
    #[clap(short, long)]
    pub cosmos_denom: String,
    /// An Ethereum private key, containing enough ETH to pay for the transaction
    #[clap(short, long, parse(try_from_str))]
    pub ethereum_key: EthPrivateKey,
    /// The address fo the Gravity contract on Ethereum
    #[clap(short, long, parse(try_from_str))]
    pub gravity_contract_address: EthAddress,
    /// The name value for the ERC20 contract, must mach Cosmos denom metadata in order to be adopted
    #[clap(short, long)]
    pub erc20_name: String,
    /// The symbol value for the ERC20 contract, must mach Cosmos denom metadata in order to be adopted
    #[clap(short, long)]
    pub erc20_symbol: String,
    /// The decimals value for the ERC20 contract, must mach Cosmos denom metadata in order to be adopted
    #[clap(short, long)]
    pub erc20_decimals: u8,
}

/// Manage keys
#[derive(Clap)]
#[clap(setting = AppSettings::ColoredHelp)]
pub struct KeyOpts {
    #[clap(subcommand)]
    subcmd: KeysSubcommand,
}

#[derive(Clap)]
pub enum KeysSubcommand {
    SetOrchestratorAddress(SetOrchestratorAddress),
}

/// Register delegate keys for the Gravity Orchestrator.
/// this is a mandatory part of setting up a Gravity Orchestrator
/// If you would like sign using a ledger see `cosmos tx gravity set-orchestrator-address` instead
#[derive(Clap)]
#[clap(setting = AppSettings::ColoredHelp)]
pub struct SetOrchestratorAddress {
    /// The Cosmos private key of the validator
    #[clap(short, long, parse(try_from_str))]
    validator_phrase: CosmosPrivateKey,
    /// (Optional) The Ethereum private key to register, will be generated if not provided
    #[clap(short, long, parse(try_from_str))]
    ethereum_key: Option<EthPrivateKey>,
    /// (Optional) The phrase for the Cosmos key to register, will be generated if not provided.
    #[clap(short, long, parse(try_from_str))]
    cosmos_phrase: Option<CosmosPrivateKey>,
    ///The prefix for Addresses on this chain (eg 'cosmos')
    #[clap(short, long)]
    address_prefix: String,
    /// (Optional) The Cosmos RPC url, usually the validator. Default is localhost:9090
    #[clap(short, long)]
    cosmos_grpc: Option<String>,
    /// The Cosmos Denom in which to pay Cosmos chain fees
    #[clap(short, long)]
    fees: String,
}
