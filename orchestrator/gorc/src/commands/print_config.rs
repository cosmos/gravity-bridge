use crate::config::GorcConfig;
use crate::{application::APP, prelude::*};
use abscissa_core::{Application, Command, Clap, Runnable};

/// Command for printing configurations
#[derive(Command, Debug, Default, Clap)]
pub struct PrintConfigCmd {
    #[clap(short, long)]
    show_default: bool,
}

impl Runnable for PrintConfigCmd {
    fn run(&self) {
        let config = if self.show_default {
            GorcConfig::default()
        } else {
            let config = APP.config();
            GorcConfig {
                keystore: config.keystore.to_owned(),
                gravity: config.gravity.to_owned(),
                ethereum: config.ethereum.to_owned(),
                cosmos: config.cosmos.to_owned(),
                metrics: config.metrics.to_owned(),
            }
        };

        print!("{}", toml::to_string(&config).unwrap());
    }
}
