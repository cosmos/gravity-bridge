use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct ShowEthKeyCmd {
    #[options(short = "n", long = "name", help = "show [name]")]
    pub name: String,
}

/// The `gorc keys eth show [name]` subcommand: show keys
impl Runnable for ShowEthKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to list keys
    }
}
