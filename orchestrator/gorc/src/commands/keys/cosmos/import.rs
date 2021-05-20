use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct ImportCosmosKeyCmd {
    #[options(
        short = "n",
        long = "name",
        help = "import private key [name] [mnemnoic]"
    )]
    pub name: String,

    #[options(
        short = "m",
        long = "mnemnoic",
        help = "import private key [name] [mnemnoic]"
    )]
    pub mnemnoic: String,
}

/// The `gork keys cosmos import [name] [mnemnoic]` subcommand: import key
impl Runnable for ImportCosmosKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to import key
    }
}
