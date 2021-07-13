//! `cosmos keys` subcommand

mod add;
mod delete;
mod import;
mod list;
mod show;
mod rename;

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum CosmosKeysCmd {
    #[options(help = "add [name]")]
    Add(add::AddCosmosKeyCmd),

    #[options(help = "import [name] [mnemnoic]")]
    Import(import::ImportCosmosKeyCmd),

    #[options(help = "delete [name]")]
    Delete(delete::DeleteCosmosKeyCmd),

    #[options(help = "rename [name] [new-name]")]
    Rename(rename::RenameCosmosKeyCmd),

    #[options(help = "list")]
    List(list::ListCosmosKeyCmd),

    #[options(help = "show [name]")]
    Show(show::ShowCosmosKeyCmd),
}

