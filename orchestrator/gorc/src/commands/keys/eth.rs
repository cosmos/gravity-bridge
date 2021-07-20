mod add;
mod delete;
mod import;
mod list;
mod recover;
mod rename;
mod show;

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum EthKeysCmd {
    #[options(help = "add [name]")]
    Add(add::AddEthKeyCmd),

    #[options(help = "delete [name]")]
    Delete(delete::DeleteEthKeyCmd),

    #[options(help = "import [name] (private-key)")]
    Import(import::ImportEthKeyCmd),

    #[options(help = "list")]
    List(list::ListEthKeyCmd),

    #[options(help = "recover [name] (bip39-mnemonic)")]
    Recover(recover::RecoverEthKeyCmd),

    #[options(help = "rename [name] [new-name]")]
    Rename(rename::RenameEthKeyCmd),

    #[options(help = "show [name]")]
    Show(show::ShowEthKeyCmd),
}
