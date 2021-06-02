//! `deploy` subcommand

use crate::{application::APP, prelude::*, utils::*};
use abscissa_core::{Command, Options, Runnable};
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use ethereum_gravity::deploy_erc20::deploy_erc20;
use gravity_proto::gravity::DenomToErc20Request;
use gravity_utils::connection_prep::{check_for_eth, check_for_fee_denom, create_rpc_connections};
use std::time::Instant;
use std::{process::exit, time::Duration, u128};
use tokio::time::sleep as delay_for;

fn lookup_eth_key(key: String) -> EthPrivateKey {
    todo!()
}

#[derive(Command, Debug, Options)]
pub enum DeployCmd {
    #[options(help = "cosmos-erc20 [denom] [erc20_name] [erc20_symbol] [erc20_decimals]")]
    Cosmos_Erc20(Cosmos_Erc20),
}

impl Runnable for DeployCmd {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Options)]
pub struct Cosmos_Erc20 {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Cosmos_Erc20 {
    /// Start the application.
    fn run(&self) {
        assert!(self.free.len() == 4);
        let denom = self.free[0].clone();
        let erc20_name = self.free[1].clone();
        let erc20_symbol = self.free[2].clone();
        let erc20_decimals = self.free[3].clone();
        let from_eth_key = self.free[4].clone();

        let config = APP.config();
        let cosmos_prefix = config.cosmos.prefix.clone();
        let eth_rpc = config.ethereum.rpc.clone();
        let cosmso_grpc = config.cosmos.grpc.clone();
        let contract_address: EthAddress = config
            .gravity
            .contract
            .clone()
            .parse()
            .expect("Expected config.gravity.contract to be an Eth ddress");

        abscissa_tokio::run(&APP, async {
            let connections =
                create_rpc_connections(cosmos_prefix, Some(cosmso_grpc), Some(eth_rpc), TIMEOUT)
                    .await;
            let web3 = connections.web3.unwrap();
            let mut grpc = connections.grpc.unwrap();

            let res = grpc
                .denom_to_erc20(DenomToErc20Request {
                    denom: denom.clone(),
                })
                .await;
            if let Ok(val) = res {
                println!(
                    "Asset {} already has ERC20 representation {}",
                    denom,
                    val.into_inner().erc20
                );
                exit(1);
            }

            let res = deploy_erc20(
                denom.clone(),
                erc20_name,
                erc20_symbol,
                erc20_decimals.parse().unwrap(),
                contract_address,
                &web3,
                Some(TIMEOUT),
                lookup_eth_key(from_eth_key),
                vec![],
            )
            .await
            .unwrap();

            let start = Instant::now();
            loop {
                let res = grpc
                    .denom_to_erc20(DenomToErc20Request {
                        denom: denom.clone(),
                    })
                    .await;

                if let Ok(val) = res {
                    println!(
                        "Asset {} has accepted new ERC20 representation {}",
                        denom,
                        val.into_inner().erc20
                    );
                    exit(0);
                }

                if Instant::now() - start > Duration::from_secs(100) {
                    println!(
                    "Your ERC20 contract was not adopted, double check the metadata and try again"
                );
                    exit(1);
                }
                delay_for(Duration::from_secs(1)).await;
            }
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}
