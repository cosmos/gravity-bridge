use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct AddCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "add private key [name]")]
    pub name: String,
}

/// The `gork keys cosmos add [name] ` subcommand: add private key
impl Runnable for AddCosmosKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to save private key to fs
    }
}