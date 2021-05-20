//! `keys` subcommand

mod cosmos;
mod eth;

use abscissa_core::{Command, Options, Runnable};

use crate::commands::keys::cosmos::CosmosKeysCmd;
use crate::commands::keys::eth::EthKeysCmd;

/// `keys` subcommand
///
/// The `Options` proc macro generates an option parser based on the struct
/// definition, and is defined in the `gumdrop` crate. See their documentation
/// for a more comprehensive example:
///
/// <https://docs.rs/gumdrop/>
#[derive(Command, Debug, Options)]
pub enum KeysCmd {
    CosmosKeysCmd(CosmosKeysCmd),
    EthKeysCmd(EthKeysCmd),
}

impl Runnable for KeysCmd {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}
