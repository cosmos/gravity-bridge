mod add;
mod delete;
mod list;
mod recover;
mod rename;
mod show;

use abscissa_core::{Command, Clap, Runnable};

#[derive(Command, Debug, Clap, Runnable)]
pub enum CosmosKeysCmd {
    Add(add::AddCosmosKeyCmd),

    Delete(delete::DeleteCosmosKeyCmd),

    Recover(recover::RecoverCosmosKeyCmd),

    Rename(rename::RenameCosmosKeyCmd),

    List(list::ListCosmosKeyCmd),

    Show(show::ShowCosmosKeyCmd),
}
