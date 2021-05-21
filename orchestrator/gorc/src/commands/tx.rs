//! `tx` subcommand

mod cosmos;

mod eth;

use abscissa_core::{Command, Options, Runnable};

/// `tx` subcommand
///
/// The `Options` proc macro generates an option parser based on the struct
/// definition, and is defined in the `gumdrop` crate. See their documentation
/// for a more comprehensive example:
///
/// <https://docs.rs/gumdrop/>
#[derive(Command, Debug, Options)]
pub enum TxCmd {
    Cosmos(cosmos::Cosmos),
    Eth(eth::Eth),
}

impl Runnable for TxCmd {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}
