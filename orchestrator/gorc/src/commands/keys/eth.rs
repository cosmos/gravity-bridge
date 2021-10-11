mod add;
mod delete;
mod import;
mod list;
mod recover;
mod rename;
mod show;

use abscissa_core::{Command, Clap, Runnable};

#[derive(Command, Debug, Clap, Runnable)]
pub enum EthKeysCmd {
    Add(add::AddEthKeyCmd),

    Delete(delete::DeleteEthKeyCmd),

    Import(import::ImportEthKeyCmd),

    List(list::ListEthKeyCmd),

    Recover(recover::RecoverEthKeyCmd),

    Rename(rename::RenameEthKeyCmd),

    Show(show::ShowEthKeyCmd),
}
