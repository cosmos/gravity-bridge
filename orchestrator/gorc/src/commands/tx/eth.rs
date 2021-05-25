//! `eth subcommands` subcommand

use abscissa_core::{Command, Options, Runnable};
use crate::{prelude::*,application::APP};


#[derive(Command, Debug, Options)]
pub enum Eth{
    SendToCosmos(SendToCosmos),
    Send(Send),
}

impl Runnable for Eth{
    fn run(&self){

    }
}

#[derive(Command, Debug, Options)]
pub struct SendToCosmos{
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,

}

impl Runnable for SendToCosmos{
    fn run(&self){

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

impl Runnable for Send{
    fn run(&self){

    }
}