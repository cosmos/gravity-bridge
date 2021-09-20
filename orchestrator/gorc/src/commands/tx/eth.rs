//! `eth subcommands` subcommand

use crate::{application::APP, prelude::*, utils::*};
use abscissa_core::{Command, Options, Runnable};
use clarity::Address as EthAddress;
use clarity::{PrivateKey as EthPrivateKey, Uint256};
use deep_space::address::Address as CosmosAddress;
use ethereum_gravity::send_to_cosmos::send_to_cosmos;
use gravity_utils::connection_prep::{check_for_eth, create_rpc_connections};

#[derive(Command, Debug, Options)]
pub enum Eth {
    #[options(
        help = "send-to-cosmos [from-eth-key][to-cosmos-addr] [erc20 conract] [erc20 amount] [[--times=int]]"
    )]
    SendToCosmos(SendToCosmos),
    #[options(help = "send [from-key] [to-addr] [amount] [token-contract]")]
    Send(Send),
}

impl Runnable for Eth {
    fn run(&self) {}
}

#[derive(Command, Debug, Options)]
pub struct SendToCosmos {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

fn lookup_eth_key(_key: String) -> EthPrivateKey {
    todo!()
}

impl Runnable for SendToCosmos {
    fn run(&self) {
        assert!(self.free.len() == 4);
        let from_eth_key = self.free[0].clone();
        let to_cosmos_addr: CosmosAddress = self.free[1]
            .clone()
            .parse()
            .expect("Expected a valid Cosmos Address");
        let erc20_contract: EthAddress = self.free[2]
            .clone()
            .parse()
            .expect("Expected a valid Eth Address");
        let erc20_amount = self.free[3].clone();
        let eth_key = lookup_eth_key(from_eth_key);

        println!(
            "Sending from Eth address {}",
            eth_key.to_public_key().unwrap()
        );
        let config = APP.config();
        let cosmos_prefix = config.cosmos.prefix.clone();
        let cosmso_grpc = config.cosmos.grpc.clone();
        let eth_rpc = config.ethereum.rpc.clone();
        let contract_address: EthAddress = config
            .gravity
            .contract
            .clone()
            .parse()
            .expect("Expected config.gravity.contract to be an Eth ddress");

        abscissa_tokio::run_with_actix(&APP, async {
            let connections =
                create_rpc_connections(cosmos_prefix, Some(cosmso_grpc), Some(eth_rpc), TIMEOUT)
                    .await;
            let web3 = connections.web3.unwrap();
            let ethereum_public_key = eth_key.to_public_key().unwrap();
            check_for_eth(ethereum_public_key, &web3).await;

            let amount: Uint256 = erc20_amount
                .parse()
                .expect("Expected amount in xx.yy format");

            let erc20_balance = web3
                .get_erc20_balance(erc20_contract, ethereum_public_key)
                .await
                .expect("Failed to get balance, check ERC20 contract address");

            if erc20_balance == 0u8.into() {
                panic!(
                    "You have zero {} tokens, please double check your sender and erc20 addresses!",
                    erc20_contract
                );
            }
            println!(
                "Sending {} / {} to Cosmos from {} to {}",
                amount, erc20_contract, ethereum_public_key, to_cosmos_addr
            );
            // we send some erc20 tokens to the gravity contract to register a deposit
            let res = send_to_cosmos(
                erc20_contract,
                contract_address,
                amount.clone(),
                to_cosmos_addr,
                eth_key,
                Some(TIMEOUT),
                &web3,
                vec![],
            )
            .await;
            match res {
                Ok(tx_id) => println!("Send to Cosmos txid: {:#066x}", tx_id),
                Err(e) => println!("Failed to send tokens! {:?}", e),
            }
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
            std::process::exit(1);
        });
    }
}

#[derive(Command, Debug, Options)]
pub struct Send {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,
}

impl Runnable for Send {
    fn run(&self) {}
}
