use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct ListEthKeyCmd {
    #[options(short = "n", long = "name", help = "list keys")]
    pub name: String,
}

/// The `gorc keys eth list` subcommand: list keys
impl Runnable for ListEthKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to list keys
    }
}
