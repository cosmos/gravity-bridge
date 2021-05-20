use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct AddEthKeyCmd {
    #[options(short = "n", long = "name", help = "add [name]")]
    pub name: String,
}

/// The `gork eth add [name] ` subcommand: add private key
impl Runnable for AddEthKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to save private key to fs
    }
}
