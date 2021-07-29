mod erc20;
use erc20::Erc20;

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum DeployCmd {
    #[options(
        name = "erc20",
        help = "deploy an ERC20 representation of a cosmos denom"
    )]
    Erc20(Erc20),
}
