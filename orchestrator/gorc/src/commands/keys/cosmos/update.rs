use super::*;
use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct UpdateCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "update [name] [new-name]")]
    pub name: String,

    #[options(short = "n", long = "name", help = "update [name] [new-name]")]
    pub new_name: String,
}

/// The `gork cosmos update [name] [new-name]` subcommand: show keys
impl Runnable for UpdateCosmosKeyCmd {
    fn run(&self) {
        /// todo(shella): glue with signatory crate to update keys
    }
}