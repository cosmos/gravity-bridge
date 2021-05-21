#[macro_use]
extern crate log;

use crate::args::{ClientSubcommand, SubCommand};
use args::Opts;
use clap::Clap;
use client::cosmos_to_eth::cosmos_to_eth;
use client::deploy_erc20_representation::deploy_erc20_representation;
use client::eth_to_cosmos::eth_to_cosmos;
use env_logger::Env;

mod args;
mod client;
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
    let address_prefix = if let Some(p) = opts.address_prefix {
        p
    } else {
        deep_space::Address::DEFAULT_PREFIX.to_string()
    };

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
        SubCommand::Keys(key_opts) => {}
        SubCommand::Orchestrator(orchestrator_opts) => {}
        SubCommand::Relayer(relayer_opts) => {}
    }

    // this may be unreachable
    panic!("No valid subcommand found!");
}
