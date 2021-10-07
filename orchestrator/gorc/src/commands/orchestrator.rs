mod start;

use abscissa_core::{Command, Clap, Runnable};

#[derive(Command, Debug, Clap, Runnable)]
pub enum OrchestratorCmd {
	#[clap(name = "start")]
	Start(start::StartCommand),
}
