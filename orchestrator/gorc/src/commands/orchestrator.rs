mod start;

use abscissa_core::{Command, Clap, Runnable};

/// Management commannds for the orchestrator
#[derive(Command, Debug, Clap, Runnable)]
pub enum OrchestratorCmd {
	Start(start::StartCommand),
}
