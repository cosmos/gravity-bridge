//! `cosmos subcommands` subcommand

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options)]
pub enum Cosmos {
    #[options(help = "balance [key-name]")]
    Balance(Balance),
    #[options(help = "gravity-keys [key-name] ")]
    Gravity_keys(Gravity_keys),
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
        let key_name = self.free[0].clone();
    }
}

#[derive(Command, Debug, Options)]
pub struct Gravity_keys {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Gravity_keys {
    /// Start the application.
    fn run(&self) {
        assert!(self.free.len() == 1);
        let key_name = self.free[0].clone();
    }
}
