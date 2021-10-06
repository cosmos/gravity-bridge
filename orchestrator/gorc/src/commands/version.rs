//! `version` subcommand

#![allow(clippy::never_loop)]

use super::GorcCmd;
use abscissa_core::{Application, Command, Clap, Runnable};

/// `version` subcommand
#[derive(Command, Debug, Default, Clap)]
pub struct VersionCmd {}

impl Runnable for VersionCmd {
    /// Print version message
    fn run(&self) {
        println!("{} {}", GorcCmd::name(), GorcCmd::version());
    }
}
