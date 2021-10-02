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
#[derive(Command, Debug, Options, Runnable)]
pub enum KeysCmd {
    #[options(name = "cosmos", help = "Cosmos key commands")]
    CosmosKeysCmd(CosmosKeysCmd),

    #[options(name = "eth", help = "Ethereum key commands.")]
    EthKeysCmd(EthKeysCmd),
}
