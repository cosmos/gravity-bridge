//! `cosmos subcommands` subcommand

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options)]
pub enum Cosmos{
    SendToEth(SendToEth),
    Send(Send),
}

impl Runnable for Cosmos {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Options)]
pub struct SendToEth{

}

impl Runnable for SendToEth {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Options)]
pub struct Send{

}

impl Runnable for Send {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}