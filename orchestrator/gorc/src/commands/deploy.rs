//! `deploy` subcommand

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options)]
pub enum DeployCmd {
    #[options(help = "cosmos-erc20 [denom] [erc20_name] [erc20_symbol] [erc20_decimals]")]
    Cosmos_Erc20(Cosmos_Erc20),
}

impl Runnable for DeployCmd {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Options)]
pub struct Cosmos_Erc20 {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Cosmos_Erc20 {
    /// Start the application.
    fn run(&self) {
        assert!(self.free.len() == 4);
        let denom = self.free[0].clone();
        let erc20_name = self.free[1].clone();
        let erc20_symbol = self.free[2].clone();
        let erc20_decimals = self.free[3].clone();
    }
}
