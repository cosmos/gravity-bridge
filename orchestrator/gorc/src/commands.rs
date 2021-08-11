//! Gorc Subcommands
//! This is where you specify the subcommands of your application.

mod deploy;
mod keys;
mod orchestrator;
mod print_config;
mod query;
mod sign_delegate_keys;
mod tests;
mod tx;
mod version;
mod cosmos_to_eth;
mod print_config;

use self::{
    keys::KeysCmd, orchestrator::OrchestratorCmd, print_config::PrintConfigCmd, query::QueryCmd,
    tests::TestsCmd, tx::TxCmd, version::VersionCmd,
};
use crate::config::GorcConfig;
use abscissa_core::{Command, Configurable, Help, Options, Runnable};
use std::path::PathBuf;

/// Gorc Configuration Filename
pub const CONFIG_FILE: &str = "gorc.toml";

/// Gorc Subcommands
#[derive(Command, Debug, Options, Runnable)]
pub enum GorcCmd {
    #[options(help = "get usage information")]
    Help(Help<Self>),

    #[options(help = "key management commands")]
    Keys(KeysCmd),

    #[options(help = "orchestrator")]
    Orchestrator(OrchestratorCmd),

    #[options(help = "print config file template")]
    PrintConfig(PrintConfigCmd),

    #[options(help = "query state on either ethereum or cosmos chains")]
    Query(QueryCmd),

    #[options(help = "sign delegate keys")]
    SignDelegateKeys(sign_delegate_keys::SignDelegateKeysCmd),

    #[options(help = "run tests against configured chains")]
    Tests(TestsCmd),

    #[options(help = "create transactions on either ethereum or cosmos chains")]
    Tx(TxCmd),

    #[options(help = "display version information")]
    Version(VersionCmd),

    #[options(help = "Send Cosmos to Ethereum")]
    CosmosToEth(cosmos_to_eth::CosmosToEthCmd),
    #[options(help = "print config file template")]
    PrintConfigCmd(PrintConfigCmd),
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
