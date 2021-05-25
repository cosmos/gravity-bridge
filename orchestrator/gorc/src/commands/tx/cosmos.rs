//! `cosmos subcommands` subcommand

use abscissa_core::{Command, Options, Runnable};
use crate::{prelude::*,application::APP};
use regex::Regex;
use clarity::Uint256;
use deep_space::address::Address as CosmosAddress;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};

#[derive(Command, Debug, Options)]
pub enum Cosmos {
    #[options(help = "send-to-eth [from-cosmos-key] [to-eth-addr] [erc20-coin] [[--times=int]]")]
    SendToEth(SendToEth),
    #[options(help = "send [from-key] [to-addr] [coin-amount]")]
    Send(Send),
}

impl Runnable for Cosmos {
    /// Start the application.
    fn run(&self) {
        // Your code goes here
    }
}

#[derive(Command, Debug, Options)]
pub struct SendToEth {
    #[options(free)]
    free: Vec<String>,

    #[options(help = "print help message")]
    help: bool,

    #[options(help = "numeber of times to sent to ethereum")]
    times: Option<u32>,
}

/// TODO revisit this for higher precision while
/// still representing the number to the user as a float
/// this takes a number like 0.37 eth and turns it into wei
/// or any erc20 with arbitrary decimals
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

fn parse_denom(s: &str) -> (String, String) {
    let re_dec_amt = r#"[[:digit:]]+(?:\.[[:digit:]]+)?|\.[[:digit:]]+"#;
    let re_dnm_string = r#"[a-zA-Z][a-zA-Z0-9/]{2,127}"#;
    let decimal_regex = Regex::new(re_dec_amt).expect("Invalid Decimal Regex");
    let denom_regex = Regex::new(re_dnm_string).expect("Invalid Denom Regex");
    let amount = decimal_regex
        .find(s)
        .expect("Could not find amount in coin string");
    let denom = denom_regex
        .find(s)
        .expect("Could not find denom in coin string");
    (amount.as_str().to_string(), denom.as_str().to_string())
}

fn get_cosmos_key(key_name: &str) -> CosmosPrivateKey {
    unimplemented!()
}

impl Runnable for SendToEth {
    fn run(&self) {
        assert!(self.free.len() == 3);
        let from_cosmos_key = self.free[0].clone();
        let to_eth_addr = self.free[1].clone(); //TODO parse this to an Eth Address
        let erc_20_coin = self.free[2].clone(); // 1231234uatom
        let (amount, denom) = parse_denom(&erc_20_coin);

        let is_cosmos_originated = !denom.starts_with("gravity");

        let amount = if is_cosmos_originated {
            fraction_to_exponent(amount.parse().unwrap(), 6)
        } else {
            fraction_to_exponent(amount.parse().unwrap(), 18)
        };

        let cosmos_key = get_cosmos_key(&from_cosmos_key);

        let cosmos_address = cosmos_key.to_address("//TODO add to config file").unwrap();

        println!("Sending from Cosmos address {}", cosmos_address);

        abscissa_tokio::run(&APP, async { unimplemented!()
         }).unwrap_or_else(|e| {
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
    /// Start the application.
    fn run(&self) {
        assert!(self.free.len() == 3);
        let from_key = self.free[0].clone();
        let to_addr = self.free[1].clone();
        let coin_amount = self.free[2].clone();
    }
}
