//! `eth keys` subcommands

mod add;
mod delete;
mod import;
mod list;
mod show;
mod rename;

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum EthKeysCmd {
    #[options(help = "add [name] (password)")]
    Add(add::AddEthKeyCmd),

    #[options(help = "import [name] (mnemonic) (password)")]
    Import(import::ImportEthKeyCmd),

    #[options(help = "delete [name]")]
    Delete(delete::DeleteEthKeyCmd),

    #[options(help = "rename [name] [new-name]")]
    Rename(rename::RenameEthKeyCmd),

    #[options(help = "list")]
    List(list::ListEthKeyCmd),

    #[options(help = "show [name]")]
    Show(show::ShowEthKeyCmd),
}

impl EthKeysCmd {
    //fn run(&self){
    //}
}
