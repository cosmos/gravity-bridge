//! `eth subcommands` subcommand

use abscissa_core::{Command, Options, Runnable};
use crate::{prelude::*,application::APP};


#[derive(Command, Debug, Options)]
pub enum Eth {
    #[options(
        help = "send-to-cosmos [from-eth-key][to-cosmos-addr] [erc20 conract] [erc20 amount] [[--times=int]]"
    )]
    SendToCosmos(SendToCosmos),
    #[options(help = "send [from-key] [to-addr] [amount] [token-contract]")]
    Send(Send),
}

impl Runnable for Eth {
    fn run(&self) {}
}

#[derive(Command, Debug, Options)]
pub struct SendToCosmos {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,

    #[options(help = "numeber of times to sent to cosmos")]
    times: Option<u32>,
}



impl Runnable for SendToCosmos {
    fn run(&self) {
        assert!(self.free.len() == 4);
        let from_eth_key = self.free[0].clone();
        let to_cosmos_addr = self.free[1].clone();
        let erc20_conract = self.free[2].clone();
        let erc20_amount = self.free[3].clone();

        abscissa_tokio::run(&APP, async { unimplemented!()
        }).unwrap_or_else(|e| {
           status_err!("executor exited with error: {}", e);
           std::process::exit(1);
       });
    }
}

#[derive(Command, Debug, Options)]
pub struct Send {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Send {
    fn run(&self) {}
}
