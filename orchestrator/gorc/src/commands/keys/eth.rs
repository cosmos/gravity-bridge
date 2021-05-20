//! `eth keys` subcommands

mod add;
mod delete;
mod import;
mod list;
mod show;
mod update;

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum EthKeysCmd {
    #[options(help = "add [name]")]
    Add(add::AddEthKeyCmd),

    #[options(help = "import [name] [privkey]")]
    Import(import::ImportEthKeyCmd),

    #[options(help = "delete [name]")]
    Delete(delete::DeleteEthKeyCmd),

    #[options(help = "update [name] [new-name]")]
    Update(update::UpdateEthKeyCmd),

    #[options(help = "list")]
    List(list::ListEthKeyCmd),

    #[options(help = "show [name]")]
    Show(show::ShowEthKeyCmd)
}

impl EthKeysCmd {
    //fn run(&self){
    //}
}
