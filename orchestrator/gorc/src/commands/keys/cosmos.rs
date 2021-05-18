//! `cosmos keys` subcommand

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum Cosmos {
    #[options(help = "add [name]")]
    Add(AddCosmosKey),

    #[options(help = "import [name] [mnemnoic]")]
    Import(ImportCosmosKey),

    #[options(help = "delete [name]")]
    Delete(DeleteCosmosKey),

    #[options(help = "update [name] [new-name]")]
    Update(UpdateCosmosKey),

    #[options(help = "list")]
    List(ListCosmosKey),

    #[options(help = "show [name]")]
    Show(ShowCosmosKey)
}


impl Runnable for Cosmos {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}