use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct RenameCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "rename [name] [new-name]")]
    pub name: String,

    #[options(help = "rename [name] [new-name]")]
    pub new_name: String,
}

/// The `gorc keys cosmos rename [name] [new-name]` subcommand: show keys
impl Runnable for RenameCosmosKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to rename keys
    }
}
