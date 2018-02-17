#[macro_use]
extern crate serde_derive;
extern crate docopt;
extern crate web3;
#[macro_use]
extern crate log;
extern crate env_logger;
extern crate tokio_core;
extern crate tokio_timer;
extern crate ethabi;
#[macro_use]
extern crate ethabi_derive;
#[macro_use]
extern crate ethabi_contract;
extern crate futures;

use web3::transports::ipc::Ipc;
use web3::types::{Address, BlockNumber, FilterBuilder, Log};
use web3::api::{self, Namespace};
use web3::Transport;
use tokio_core::reactor::Core;
use futures::Future;
use ethabi::RawLog;

// makes the contract available as toy::Toy
use_contract!(toy, "Toy", "SolidityToy.abi");

const USAGE: &'static str = "
Usage: sfeature/eth_witnessigner [--contract=<address>] [--ipc=<path.ipc>]

Options:
    --ipc=<path>                Path to unix socket. [default: /Users/adrianbrink/.peggy/jsonrpc.ipc]
    --contract=<address>        Contract address.    [default: 0xdd1cB580B505b59962Ef7a31d21CEE7234225C29]
";

#[derive(Deserialize)]
struct Args {
    flag_ipc: String,
    flag_contract: String,
}

fn sign_and_forward(log: Log) -> bool {
    true
}

fn main() {
    let args: Args = docopt::Docopt::new(USAGE)
        .and_then(|d| d.argv(std::env::args().into_iter()).deserialize())
        .unwrap_or_else(|e| e.exit());

    env_logger::init();

    // TODO store in db as ack is received from ABCI
    let mut last_block: u64 = 0;
    let mut event_loop = Core::new().unwrap();

    info!("making ipc event loop");
    let ipc = Ipc::with_event_loop(&*args.flag_ipc, &event_loop.handle())
        .expect("should be able to connect to local unix socket");

    let address: Address = args.flag_contract.parse().expect(
        "should be able to parse address",
    );

    let filter_builder = FilterBuilder::default()
        .from_block(BlockNumber::Number(0))
        .to_block(BlockNumber::Latest)
        // .limit(1)
        .address(vec![address]);

    info!("creating transport");
    let transport = api::Eth::new(&ipc);

    loop {
        let filter = filter_builder
            .clone()
            .from_block(BlockNumber::Number(last_block))
            .build();

        trace!("querying logs with filter {:?}", filter);

        let logs_fut = transport.logs(&filter);
        let logs = event_loop.run(logs_fut).unwrap();

        for log in logs {
            let block = log.block_number;
            info!("got log {:?}", block);
            let success = sign_and_forward(log);
            if success {
                last_block = block.unwrap().low_u64() + 1;
            }
        }
    }
}
