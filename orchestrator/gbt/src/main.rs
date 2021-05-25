#[macro_use]
extern crate log;
#[macro_use]
extern crate serde_derive;

use crate::config::init_config;
use crate::{
    args::{ClientSubcommand, KeysSubcommand, SubCommand},
    config::load_config,
};
use crate::{orchestrator::orchestrator, relayer::relayer};
use args::Opts;
use clap::Clap;
use client::cosmos_to_eth::cosmos_to_eth;
use client::deploy_erc20_representation::deploy_erc20_representation;
use client::eth_to_cosmos::eth_to_cosmos;
use env_logger::Env;
use keys::set_orchestrator_address::set_orchestrator_address;
use std::{path::PathBuf, process::exit};

mod args;
mod client;
mod config;
mod keys;
mod orchestrator;
mod relayer;
mod utils;

#[actix_rt::main]
async fn main() {
    env_logger::Builder::from_env(Env::default().default_filter_or("info")).init();
    // On Linux static builds we need to probe ssl certs path to be able to
    // do TLS stuff.
    openssl_probe::init_ssl_cert_env_vars();
    // parse the arguments
    let opts: Opts = Opts::parse();

    // handle global config here
    let address_prefix = opts.address_prefix;
    let home_dir: PathBuf = match (dirs::home_dir(), opts.home) {
        (_, Some(user_home)) => PathBuf::from(&user_home),
        (Some(default_home_dir), None) => default_home_dir,
        (None, None) => {
            error!("Failed to automatically determine your home directory, please provide a path to the --home argument!");
            exit(1);
        }
    };
    let _config = load_config(home_dir.clone());

    // control flow for the command structure
    match opts.subcmd {
        SubCommand::Client(client_opts) => match client_opts.subcmd {
            ClientSubcommand::EthToCosmos(eth_to_cosmos_opts) => {
                eth_to_cosmos(eth_to_cosmos_opts, address_prefix).await
            }
            ClientSubcommand::CosmosToEth(cosmos_to_eth_opts) => {
                cosmos_to_eth(cosmos_to_eth_opts, address_prefix).await
            }
            ClientSubcommand::DeployErc20Representation(deploy_erc20_opts) => {
                deploy_erc20_representation(deploy_erc20_opts, address_prefix).await
            }
        },
        SubCommand::Keys(key_opts) => match key_opts.subcmd {
            KeysSubcommand::SetOrchestratorAddress(set_orchestrator_address_opts) => {
                set_orchestrator_address(set_orchestrator_address_opts, address_prefix).await
            }
        },
        SubCommand::Orchestrator(orchestrator_opts) => {
            orchestrator(orchestrator_opts, address_prefix).await
        }
        SubCommand::Relayer(relayer_opts) => relayer(relayer_opts, address_prefix).await,
        SubCommand::Init(init_opts) => init_config(init_opts, home_dir),
    }
}
