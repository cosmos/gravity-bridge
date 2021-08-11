use crate::application::APP;
use abscissa_core::{status_err, Application, Command, Options, Runnable};
use clarity::Address as EthAddress;
use clarity::Uint256;
use deep_space::address::Address as CosmosAddress;
use ethereum_gravity::send_to_cosmos::send_to_cosmos;
use gravity_utils::connection_prep::{check_for_eth, create_rpc_connections};
use std::time::Duration;

const TIMEOUT: Duration = Duration::from_secs(60);

#[derive(Command, Debug, Default, Options)]
pub struct EthToCosmosCmd {
    #[options(
        free,
        help = "eth-to-cosmos [erc20_address] [ethereum_key] [contract_address] [cosmos_dest] [amount] [times]"
    )]
    pub args: Vec<String>,
}

impl Runnable for EthToCosmosCmd {
    fn run(&self) {
        let config = APP.config();
        let erc20_address = self.args.get(0).expect("erc20 address is required");
        let erc20_address: EthAddress = erc20_address
            .parse()
            .expect("Invalid ERC20 contract address!");

        let ethereum_key = self.args.get(1).expect("key is required");
        let ethereum_key = config.load_clarity_key(ethereum_key.clone());

        let contract_address = self.args.get(2).expect("contract address is required");
        let contract_address: EthAddress =
            contract_address.parse().expect("Invalid contract address!");

        let cosmos_prefix = config.cosmos.prefix.trim();
        let eth_rpc = config.ethereum.rpc.trim();
        abscissa_tokio::run_with_actix(&APP, async {
            let connections = create_rpc_connections(
                cosmos_prefix.to_string(),
                None,
                Some(eth_rpc.to_string()),
                TIMEOUT,
            )
            .await;
            let web3 = connections.web3.unwrap();
            let cosmos_dest = self.args.get(3).expect("cosmos destination is required");
            let cosmos_dest: CosmosAddress = cosmos_dest.parse().unwrap();
            let ethereum_public_key = ethereum_key.to_public_key().unwrap();
            check_for_eth(ethereum_public_key, &web3).await;

            let init_amount = self.args.get(4).expect("amount is required");
            let amount: Uint256 = init_amount.parse().unwrap();

            let erc20_balance = web3
                .get_erc20_balance(erc20_address, ethereum_public_key)
                .await
                .expect("Failed to get balance, check ERC20 contract address");

            let times = self.args.get(5).expect("times is required");
            let times = times.parse::<usize>().expect("cannot parse times");

            if erc20_balance == 0u8.into() {
                panic!(
                    "You have zero {} tokens, please double check your sender and erc20 addresses!",
                    contract_address
                );
            } else if amount.clone() * times.into() > erc20_balance {
                panic!(
                    "Insufficient balance {} > {}",
                    amount * times.into(),
                    erc20_balance
                );
            }

            for _ in 0..times {
                println!(
                    "Sending {} / {} to Cosmos from {} to {}",
                    init_amount.parse::<f64>().unwrap(),
                    erc20_address,
                    ethereum_public_key,
                    cosmos_dest
                );
                // we send some erc20 tokens to the gravity contract to register a deposit
                let res = send_to_cosmos(
                    erc20_address,
                    contract_address,
                    amount.clone(),
                    cosmos_dest,
                    ethereum_key,
                    Some(TIMEOUT),
                    &web3,
                    vec![],
                )
                .await;
                match res {
                    Ok(tx_id) => println!("Send to Cosmos txid: {:#066x}", tx_id),
                    Err(e) => println!("Failed to send tokens! {:?}", e),
                }
            }
        })
        .unwrap_or_else(|e| {
            status_err!("executor exited with error: {}", e);
        });
    }
}
