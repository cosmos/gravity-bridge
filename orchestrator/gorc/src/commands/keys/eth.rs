//! `eth keys` subcommands

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum EthKeysCmd {
    #[options(help = "add [name]")]
    Add(AddEthKeyCmd),

    #[options(help = "import [name] [privkey]")]
    Import(ImportEthKeyCmd),

    #[options(help = "delete [name]")]
    Delete(DeleteEthKeyCmd),

    #[options(help = "update [name] [new-name]")]
    Update(UpdateEthKeyCmd),

    #[options(help = "list")]
    List(ListEthKeyCmd),

    #[options(help = "show [name]")]
    Show(ShowEthKeyCmd)
}

impl Runnable for Eth{
    fn run(&self){

    }
}
