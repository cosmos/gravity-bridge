use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct DeleteEthKeyCmd {
    #[options(short = "n", long = "name", help = "delete key [name]")]
    pub name: String,
}

/// The `gork eth delete [name] ` subcommand: delete private key
impl Runnable for DeleteEthKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to rm key from fs
    }
}
