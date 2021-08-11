mod start;

use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Options, Runnable)]
pub enum OrchestratorCmd {
	#[options(help = "start the orchestrator")]
	Start(start::StartCommand),
}
