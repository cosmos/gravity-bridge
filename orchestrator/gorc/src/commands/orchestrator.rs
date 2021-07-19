mod start;

use abscissa_core::{Command, Options, Runnable};

/// `orchestator` subcommand
///
/// The `Options` proc macro generates an option parser based on the struct
/// definition, and is defined in the `gumdrop` crate. See their documentation
/// for a more comprehensive example:
///
/// <https://docs.rs/gumdrop/>
#[derive(Command, Debug, Options, Runnable)]
pub enum OrchestratorCmd {
	#[options(name = "start")]
	Start(start::StartCommand),
}