mod erc20;
use erc20::Erc20;

use abscissa_core::{Command, Clap, Runnable};

#[derive(Command, Debug, Clap, Runnable)]
pub enum DeployCmd {
    #[clap(name = "erc20")]
    Erc20(Erc20),
}
