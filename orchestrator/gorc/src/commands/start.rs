//! `start` subcommand - example of how to write a subcommand

use crate::{application::APP, prelude::*};
/// App-local prelude includes `app_reader()`/`app_writer()`/`app_config()`
/// accessors along with logging macros. Customize as you see fit.
use abscissa_core::{Command, Options, Runnable};

/// `start` subcommand

#[derive(Command, Debug, Options)]
pub enum StartCmd {
    /// To whom are we saying hello?
    #[options(help = "orchestrator [contract-address] [fee-denom]")]
    Orchestrator(Orchestrator),

    #[options(help = "relayer")]
    Relayer(Relayer),
}

impl Runnable for StartCmd {
    /// Start the application.
    fn run(&self) {
        //Your code goes here
    }
}

#[derive(Command, Debug, Options)]
pub struct Orchestrator {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Orchestrator {
    fn run(&self) {
        assert!(self.free.len() == 2);
        let contract_address = self.free[0].clone();
        let fee_denom = self.free[1].clone();

        abscissa_tokio::run(&APP, async { unimplemented!() }).unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}

#[derive(Command, Debug, Options)]
pub struct Relayer {
    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Relayer {
    /// Start the application.
    fn run(&self) {}
}
