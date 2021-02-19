use crate::main_loop::relayer_main_loop;
use crate::main_loop::LOOP_SPEED;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use docopt::Docopt;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use url::Url;
use web30::client::Web3;

pub mod batch_relaying;
pub mod find_latest_valset;
pub mod logic_call_relaying;
pub mod main_loop;
pub mod valset_relaying;

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate lazy_static;
#[macro_use]
extern crate log;

#[derive(Debug, Deserialize)]
struct Args {
    flag_ethereum_key: String,
    flag_cosmos_legacy_rpc: String,
    flag_cosmos_grpc: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage: {} --ethereum-key=<key> --cosmos-legacy-rpc=<url> --cosmos-grpc=<url> --ethereum-rpc=<url> --fees=<denom> --contract-address=<addr>
        Options:
            -h --help                    Show this screen.
            --ethereum-key=<ekey>        An Ethereum private key containing non-trivial funds
            --cosmos-legacy-rpc=<curl>   The Cosmos RPC url
            --cosmos-grpc=<gurl>         The Cosmos gRPC url
            --ethereum-rpc=<eurl>        The Ethereum RPC url, Geth light clients work and sync fast
            --contract-address=<addr>    The Ethereum contract address for Peggy
        About:
            The Peggy relayer component, responsible for relaying data from the Cosmos blockchain
            to the Ethereum blockchain, cosmos key and fees are optional since they are only used
            to request the creation of batches or validator sets to relay.
            for Althea-Peggy.
            Written By: {}
            Version {}",
            env!("CARGO_PKG_NAME"),
            env!("CARGO_PKG_AUTHORS"),
            env!("CARGO_PKG_VERSION"),
        );
}

#[actix_rt::main]
async fn main() {
    env_logger::init();
    // On Linux static builds we need to probe ssl certs path to be able to
    // do TLS stuff.
    openssl_probe::init_ssl_cert_env_vars();

    let args: Args = Docopt::new(USAGE.as_str())
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());
    let ethereum_key: EthPrivateKey = args
        .flag_ethereum_key
        .parse()
        .expect("Invalid Ethereum private key!");
    let peggy_contract_address: EthAddress = args
        .flag_contract_address
        .parse()
        .expect("Invalid contract address!");

    let _ = Url::parse(&args.flag_cosmos_grpc).expect("Invalid Cosmos gRPC url");
    let cosmos_grpc_url = args.flag_cosmos_grpc.trim_end_matches('/').to_string();

    let _ = Url::parse(&args.flag_ethereum_rpc).expect("Invalid Ethereum RPC url");
    let eth_url = args.flag_ethereum_rpc.trim_end_matches('/');

    let grpc_client = PeggyQueryClient::connect(cosmos_grpc_url).await.unwrap();
    let web3 = Web3::new(&eth_url, LOOP_SPEED);

    let public_eth_key = ethereum_key
        .to_public_key()
        .expect("Invalid Ethereum Private Key!");
    info!("Starting Peggy Relayer");
    info!("Ethereum Address: {}", public_eth_key);

    relayer_main_loop(ethereum_key, web3, grpc_client, peggy_contract_address).await
}
