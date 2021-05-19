use super::*;
use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct ListCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "import private key [name] [mnemnoic]")]
    pub name: String,
}

/// The `gork cosmos list` subcommand: list keys
impl Runnable for ListCosmosKeyCmd {
    fn run(&self) {
        /// todo(shella): glue with signatory crate to list keys
    }
}