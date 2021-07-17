//! `eth subcommands` subcommand

use crate::{application::APP, prelude::*};
use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options)]
pub enum Eth {
    #[options(help = "balance [key-name]")]
    Balance(Balance),

    #[options(help = "contract")]
    Contract(Contract),
}

impl Runnable for Eth {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Options)]
pub struct Balance {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Balance {
    fn run(&self) {
        assert!(self.free.len() == 1);
        let _key_name = self.free[0].clone();

        abscissa_tokio::run(&APP, async { unimplemented!() }).unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}

#[derive(Command, Debug, Options)]
pub struct Contract {
    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Contract {
    /// Start the application.
    fn run(&self) {}
}
