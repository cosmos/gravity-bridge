//! `eth keys` subcommands

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum Eth {
    #[options(help = "add [name]")]
    Add(AddEthKey),

    #[options(help = "import [name] [privkey]")]
    Import(ImportEthKey),

    #[options(help = "delete [name]")]
    Delete(DeleteEthKey),

    #[options(help = "update [name] [new-name]")]
    Update(UpdateEthKey),

    #[options(help = "list")]
    List(ListEthKey),

    #[options(help = "show [name]")]
    Show(ShowEthKey)
}

impl Runnable for Eth{
    fn run(&self){

    }
}
