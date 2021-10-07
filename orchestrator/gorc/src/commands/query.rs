//! `query` subcommand

mod cosmos;

mod eth;

use abscissa_core::{Command, Clap, Runnable};

/// `query` subcommand
///
/// The `Options` proc macro generates an option parser based on the struct
/// definition, and is defined in the `gumdrop` crate. See their documentation
/// for a more comprehensive example:
///
/// <https://docs.rs/gumdrop/>
#[derive(Command, Debug, Clap)]
pub enum QueryCmd {
    #[clap(subcommand)]
    Cosmos(cosmos::Cosmos),

    #[clap(subcommand)]
    Eth(eth::Eth),
    // Example `--foobar` (with short `-f` argument)
    // #[options(short = "f", help = "foobar path"]
    // foobar: Option<PathBuf>

    // Example `--baz` argument with no short version
    // #[options(no_short, help = "baz path")]
    // baz: Option<PathBuf>

    // "free" arguments don't have an associated flag
    // #[options(free)]
    // free_args: Vec<String>,
}

impl Runnable for QueryCmd {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}
