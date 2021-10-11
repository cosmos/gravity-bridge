mod cosmos;
mod eth;

use abscissa_core::{Command, Clap, Runnable};

use crate::commands::keys::cosmos::CosmosKeysCmd;
use crate::commands::keys::eth::EthKeysCmd;

/// Key management commands for Ethereum and Cosmos

#[derive(Command, Debug, Clap, Runnable)]
pub enum KeysCmd {
    #[clap(subcommand)]
    Cosmos(CosmosKeysCmd),

    #[clap(subcommand)]
    Eth(EthKeysCmd),
}
