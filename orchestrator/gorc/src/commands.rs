//! Gorc Subcommands
//! This is where you specify the subcommands of your application.

mod cosmos_to_eth;
mod deploy;
mod eth_to_cosmos;
mod keys;
mod orchestrator;
mod print_config;
mod query;
mod sign_delegate_keys;
mod tests;
mod tx;
mod version;

use crate::config::GorcConfig;
use abscissa_core::{Application, Command, Clap, Runnable, Configurable};
use std::path::PathBuf;

/// Gorc Configuration Filename
pub const CONFIG_FILE: &str = "gorc.toml";

/// Gorc Subcommands
#[derive(Command, Debug, Clap)]
pub enum GorcCmd {
    #[clap(short, long)]
    CosmosToEth(cosmos_to_eth::CosmosToEthCmd),

    #[clap(short, long)]
    Deploy(deploy::DeployCmd),

    #[clap(short, long)]
    EthToCosmos(eth_to_cosmos::EthToCosmosCmd),

    #[clap(short, long)]
    Keys(keys::KeysCmd),

    #[clap(short, long)]
    Orchestrator(orchestrator::OrchestratorCmd),

    #[clap(short, long)]
    PrintConfig(print_config::PrintConfigCmd),

    #[clap(short, long)]
    Query(query::QueryCmd),

    #[clap(short, long)]
    SignDelegateKeys(sign_delegate_keys::SignDelegateKeysCmd),

    #[clap(short, long)]
    Tests(tests::TestsCmd),

    #[clap(short, long)]
    Tx(tx::TxCmd),

    #[clap(short, long)]
    Version(version::VersionCmd),
}

/// This trait allows you to define how application configuration is loaded.
impl Configurable<GorcConfig> for GorcCmd {
    /// Location of the configuration file
    fn config_path(&self) -> Option<PathBuf> {
        // Check if the config file exists, and if it does not, ignore it.
        // If you'd like for a missing configuration file to be a hard error
        // instead, always return `Some(CONFIG_FILE)` here.
        let filename = PathBuf::from(CONFIG_FILE);

        if filename.exists() {
            Some(filename)
        } else {
            None
        }
    }
}
