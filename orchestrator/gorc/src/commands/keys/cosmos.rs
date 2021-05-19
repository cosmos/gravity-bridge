//! `cosmos keys` subcommand

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum CosmosKeysCmd {
    #[options(help = "add [name]")]
    Add(AddCosmosKeyCmd),

    #[options(help = "import [name] [mnemnoic]")]
    Import(ImportCosmosKeyCmd),

    #[options(help = "delete [name]")]
    Delete(DeleteCosmosKeyCmd),

    #[options(help = "update [name] [new-name]")]
    Update(UpdateCosmosKeyCmd),

    #[options(help = "list")]
    List(ListCosmosKeyCmd),

    #[options(help = "show [name]")]
    Show(ShowCosmosKeyCmd)
}


impl Runnable for Cosmos {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}