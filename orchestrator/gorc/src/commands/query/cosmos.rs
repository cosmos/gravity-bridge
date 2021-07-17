//! `cosmos subcommands` subcommand

use crate::{application::APP, prelude::*};
use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options)]
pub enum Cosmos {
    #[options(help = "balance [key-name]")]
    Balance(Balance),
    #[options(help = "gravity-keys [key-name] ")]
    GravityKeys(GravityKeys),
}

impl Runnable for Cosmos {
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
    }
}

#[derive(Command, Debug, Options)]
pub struct GravityKeys {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for GravityKeys {
    /// Start the application.
    fn run(&self) {
        assert!(self.free.len() == 1);
        let _key_name = self.free[0].clone();

        abscissa_tokio::run(&APP, async { unimplemented!() }).unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}
