use crate::application::APP;
use absicssa_core::{Application, Command, Options, Runnable};
use clarity::Address as EthAddress;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use gravity_proto::gravity::DenomToErc20Request;

#[derive(Command, Debug, Default, Options)]
pub struct CosmosToEthCmd {
    #[options(free, help = "add [name]")]
    pub args: Vec<String>,
}

pub fn fraction_to_exponent(num: f64, exponent: u8) -> Uint256 {
    let mut res = num;
    // in order to avoid floating point rounding issues we
    // multiply only by 10 each time. this reduces the rounding
    // errors enough to be ignored
    for _ in 0..exponent {
        res *= 10f64
    }
    (res as u128).into()
}

impl Runnable for CosmosToEthCmd {
    fn run(&self) {
        let gravity_denom = self.args.get(0).expect("name is required");
        let is_cosmos_originated = !gravity_denom.starts_with("gravity");
        let amount = if is_cosmos_originated {
            fraction_to_exponent(self.args.get(1).expect("name is required"), 6)
        } else {
            fraction_to_exponent(self.args.get(1).expect("name is required"), 18)
        };
        let cosmos_key =
            CosmosPrivateKey::from_phrase(&self.args.get(2).expect("name is required"), "")
                .expect("Failed to parse cosmos key phrase, does it have a password?");
        let cosmos_address = cosmos_key.to_address(&self.args.get(3).expect("name is required"));

        println!("Sending from Cosmos address {}", cosmos_address);
        let connections = create_rpc_connections(
            &self.args.get(3).expect("name is required"),
            Some(self.args.get(4).expect("name is required")),
            None,
            TIMEOUT,
        )
        .await;
        let contact = connections.contact.unwrap();
        let mut grpc = connections.grpc.unwrap();
        let res = grpc
            .denom_to_erc20(DenomToErc20Request {
                denom: gravity_denom.clone(),
            })
            .await;
        match res {
            Ok(val) => println!(
                "Asset {} has ERC20 representation {}",
                gravity_denom,
                val.into_inner().erc20
            ),
            Err(_e) => {
                println!(
                    "Asset {} has no ERC20 representation, you may need to deploy an ERC20 for it!",
                    gravity_denom
                );
                exit(1);
            }
        }
        let amount = Coin {
            amount,
            denom: gravity_denom.clone(),
        };
        let bridge_fee = Coin {
            denom: gravity_denom.clone(),
            amount: 1u64.into(),
        };

        let eth_dest: EthAddress = self.args.get(5).parse().expect("name is required");
        check_for_fee_denom(&gravity_denom, cosmos_address, &contact).await;

        let balances = contact
            .get_balances(cosmos_address)
            .await
            .expect("Failed to get balances!");
        let mut found = None;
        for coin in balances.iter() {
            if coin.denom == gravity_denom {
                found = Some(coin);
            }
        }

        println!("Cosmos balances {:?}", balances);

        if found.is_none() {
            panic!("You don't have any {} tokens!", gravity_denom);
        } else if amount.amount.clone() * times.into() >= found.clone().unwrap().amount
            && times == 1
        {
            if is_cosmos_originated {
                panic!("Your transfer of {} {} tokens is greater than your balance of {} tokens. Remember you need some to pay for fees!", print_atom(amount.amount), gravity_denom, print_atom(found.unwrap().amount.clone()));
            } else {
                panic!("Your transfer of {} {} tokens is greater than your balance of {} tokens. Remember you need some to pay for fees!", print_eth(amount.amount), gravity_denom, print_eth(found.unwrap().amount.clone()));
            }
        } else if amount.amount.clone() * times.into() >= found.clone().unwrap().amount {
            if is_cosmos_originated {
                panic!("Your transfer of {} * {} {} tokens is greater than your balance of {} tokens. Try to reduce the amount or the --times parameter", print_atom(amount.amount), times, gravity_denom, print_atom(found.unwrap().amount.clone()));
            } else {
                panic!("Your transfer of {} * {} {} tokens is greater than your balance of {} tokens. Try to reduce the amount or the --times parameter", print_eth(amount.amount), times, gravity_denom, print_eth(found.unwrap().amount.clone()));
            }
        }

        for _ in 0..times {
            println!(
                "Locking {} / {} into the batch pool",
                args.flag_amount.unwrap(),
                gravity_denom
            );
            let res = send_to_eth(
                cosmos_key,
                eth_dest,
                amount.clone(),
                bridge_fee.clone(),
                &contact,
            )
            .await;
            match res {
                Ok(tx_id) => println!("Send to Eth txid {}", tx_id.txhash),
                Err(e) => println!("Failed to send tokens! {:?}", e),
            }
        }

        if !args.flag_no_batch {
            println!("Requesting a batch to push transaction along immediately");
            send_request_batch(cosmos_key, gravity_denom, bridge_fee, &contact)
                .await
                .expect("Failed to request batch");
        } else {
            println!("--no-batch specified, your transfer will wait until someone requests a batch for this token type")
        }
    }
}
