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
    #[clap(name = "add")]
    Add(add::AddEthKeyCmd),

    #[clap(name = "delete")]
    Delete(delete::DeleteEthKeyCmd),

    #[clap(name = "import")]
    Import(import::ImportEthKeyCmd),

    #[clap(name = "list")]
    List(list::ListEthKeyCmd),

    #[clap(name = "recover")]
    Recover(recover::RecoverEthKeyCmd),

    #[clap(name = "rename")]
    Rename(rename::RenameEthKeyCmd),

    #[clap(name = "show")]
    Show(show::ShowEthKeyCmd),
}
