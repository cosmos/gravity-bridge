use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct ShowCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "show [name]")]
    pub name: String,
}

/// The `gorc keys cosmos show [name]` subcommand: show keys
impl Runnable for ShowCosmosKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to list keys
    }
}
