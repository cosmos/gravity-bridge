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
}

impl Runnable for SendToCosmos{
    fn run(&self){

    }
}

#[derive(Command, Debug, Options)]
pub struct Send {

}

impl Runnable for Send{
    fn run(&self){

    }
}