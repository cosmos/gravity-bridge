use crate::{application::APP, prelude::*};
use abscissa_core::{Command, Options, Runnable};
use ethereum_gravity::deploy_erc20::deploy_erc20;
use gravity_proto::gravity::{DenomToErc20ParamsRequest, DenomToErc20Request};
use gravity_utils::connection_prep::{check_for_eth, create_rpc_connections};
use std::convert::TryFrom;
use std::process::exit;
use std::time::{Duration, Instant};
use tokio::time::sleep as delay_for;

#[derive(Command, Debug, Options)]
pub struct Erc20 {
    #[options(free, help = "denom")]
    args: Vec<String>,

    #[options(help = "ethereum key name")]
    ethereum_key: String,
}

impl Runnable for Erc20 {
    fn run(&self) {
        abscissa_tokio::run_with_actix(&APP, async {
            self.deploy().await;
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            exit(1);
        });
    }
}

impl Erc20 {
    async fn deploy(&self) {
        let denom = self.args.get(0).expect("denom is required");

        let config = APP.config();

        let contract_address = config
            .gravity
            .contract
            .parse()
            .expect("Could not parse gravity contract address");

        let timeout = Duration::from_secs(500);
        let connections = create_rpc_connections(
            config.cosmos.prefix.clone(),
            Some(config.cosmos.grpc.clone()),
            Some(config.ethereum.rpc.clone()),
            timeout,
        )
        .await;

        let mut grpc = connections.grpc.clone().unwrap();
        let web3 = connections.web3.clone().unwrap();

        let ethereum_key = config.load_clarity_key(self.ethereum_key.clone());
        let ethereum_public_key = ethereum_key.to_public_key().unwrap();
        check_for_eth(ethereum_public_key, &web3).await;

        let req = DenomToErc20ParamsRequest {
            denom: denom.clone(),
        };

        let res = grpc
            .denom_to_erc20_params(req)
            .await
            .expect("Couldn't get erc-20 params")
            .into_inner();

        println!("Starting deploy of ERC20");

        let res = deploy_erc20(
            res.base_denom,
            res.erc20_name,
            res.erc20_symbol,
            u8::try_from(res.erc20_decimals).unwrap(),
            contract_address,
            &web3,
            Some(timeout),
            ethereum_key,
            vec![],
        )
        .await
        .expect("Could not deploy ERC20");

        println!("We have deployed ERC20 contract {:#066x}, waiting to see if the Cosmos chain choses to adopt it", res);

        let start = Instant::now();
        loop {
            let req = DenomToErc20Request {
                denom: denom.clone(),
            };

            let res = grpc.denom_to_erc20(req).await;

            if let Ok(val) = res {
                let val = val.into_inner();
                println!(
                    "Asset {} has accepted new ERC20 representation {}",
                    denom, val.erc20
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
    }
}
