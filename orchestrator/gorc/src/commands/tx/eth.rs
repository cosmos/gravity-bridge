//! `eth subcommands` subcommand

use abscissa_core::{Command, Options, Runnable};

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