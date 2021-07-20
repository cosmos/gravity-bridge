mod add;
mod delete;
mod list;
mod recover;
mod rename;
mod show;

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum CosmosKeysCmd {
    #[options(help = "add [name]")]
    Add(add::AddCosmosKeyCmd),

    #[options(help = "delete [name]")]
    Delete(delete::DeleteCosmosKeyCmd),

    #[options(help = "import [name] (bip39-mnemnoic)")]
    Recover(recover::RecoverCosmosKeyCmd),

    #[options(help = "rename [name] [new-name]")]
    Rename(rename::RenameCosmosKeyCmd),

    #[options(help = "list")]
    List(list::ListCosmosKeyCmd),

    #[options(help = "show [name]")]
    Show(show::ShowCosmosKeyCmd),
}
