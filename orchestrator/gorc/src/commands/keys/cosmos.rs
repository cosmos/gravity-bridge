mod add;
mod delete;
mod list;
mod recover;
mod rename;
mod show;

use abscissa_core::{Command, Clap, Runnable};

#[derive(Command, Debug, Clap, Runnable)]
pub enum CosmosKeysCmd {
    #[clap(name = "add")]
    Add(add::AddCosmosKeyCmd),

    #[clap(name = "delete")]
    Delete(delete::DeleteCosmosKeyCmd),

    #[clap(name = "import")]
    Recover(recover::RecoverCosmosKeyCmd),

    #[clap(name = "rename")]
    Rename(rename::RenameCosmosKeyCmd),

    #[clap(name = "list")]
    List(list::ListCosmosKeyCmd),

    #[clap(name = "show")]
    Show(show::ShowCosmosKeyCmd),
}
