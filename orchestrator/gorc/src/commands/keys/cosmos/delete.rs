
use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct DeleteCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "delete key [name]")]
    pub name: String,
}

/// The `gork keys cosmos add [name] ` subcommand: delete private key
impl Runnable for DeleteCosmosKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to rm key from fs
    }
}