//! `tests` subcommand

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options)]
pub enum TestsCmd {
    #[options(help = "runner")]
    Runner(Runner),
}

impl Runnable for TestsCmd {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Options)]
pub struct Runner {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Runner {
    fn run(&self) {}
}
