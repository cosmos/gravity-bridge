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
#[macro_use]
extern crate error_chain;
extern crate ed25519_dalek;
extern crate protobuf;
extern crate sha2;

mod errors;
mod tx;

use web3::transports::ipc::Ipc;
use web3::types::{Address, BlockNumber, FilterBuilder, Log, Bytes};
use web3::api::{self, Namespace};
use web3::Transport;
use tokio_core::reactor::Core;
use futures::Future;
use ethabi::RawLog;
use errors::Result;
use peggy::logs::Lock;
use std::{thread, time};
use tx::{WitnessTx, LockMsg};
use ed25519_dalek::{Keypair, Signature, SECRET_KEY_LENGTH, SIGNATURE_LENGTH};
use protobuf::Message;
use std::iter::Iterator;
use sha2::Sha512;

// makes the contract available as toy::Toy
use_contract!(peggy, "Peggy", "Peggy.abi");

const USAGE: &'static str = "
Usage: feature/eth_witnessigner [--contract=<address>] [--ipc=<path.ipc>]

Options:
    --ipc=<path>                Path to unix socket. [default: /Users/adrianbrink/.peggy/jsonrpc.ipc]
    --contract=<address>        Contract address.    [default: 0xdd1cB580B505b59962Ef7a31d21CEE7234225C29]
";

// $HOME/.local/share/io.parity.ethereum/jsonrpc.ipc
// 0x2712a785ac11528e0b3650e3aaae2ede1508c649

#[derive(Deserialize)]
struct Args {
    flag_ipc: String,
    flag_contract: String,
}

enum WitnessLog {
    Lock(Lock)
}

impl From<Lock> for WitnessLog {
    fn from(item: Lock) -> Self {
        WitnessLog::Lock(item)
    }
}


fn new_witness(rawlog: RawLog, peggy: &peggy::Peggy) -> WitnessLog {
    let lock = peggy.events().lock().parse_log(rawlog).map(|x| WitnessLog::from(x));
    Err(()).or(lock).expect("New witness")
}

fn sign_and_wrap_lock(log: Lock, keypair: Keypair, sequence: i64) -> WitnessTx {
    let msg = LockMsg::new();
    msg.set_dest(log.to);
    msg.set_value(log.value.as_u64());
    msg.set_token(log.token.to_vec());

    let signbytes = msg.clone().write_to_bytes().unwrap();
    let signature = keypair.sign::<Sha512>(&signbytes).to_bytes();

    let tx = WitnessTx::new();
    tx.set_lock(msg);
    tx.set_signature(signature.to_vec());
    tx.set_sequence(sequence);

    tx
}

fn sign_and_wrap(log: WitnessLog, keypair: Keypair, sequence: i64) -> WitnessTx {
    match log {
        WitnessLog::Lock(l) => sign_and_wrap_lock(l, keypair, sequence)
    }
}

fn gen_rawlog(log: Log) -> RawLog {
    RawLog {
        topics: log.topics.into_iter().map(|t| From::from(t.0)).collect(),
        data: log.data.0,
    }
}



fn sign_and_forward(tx: WitnessTx) -> bool {
    true
}

fn main() {
    let args: Args = docopt::Docopt::new(USAGE)
        .and_then(|d| d.argv(std::env::args().into_iter()).deserialize())
        .unwrap_or_else(|e| e.exit());

    env_logger::init();

    // TODO store in db as ack is received from ABCI
    let mut event_loop = Core::new().unwrap();
    let delay = 2; // for testing purpose.
    let mut last_block: u64 = 0;

    println!("making ipc event loop");
    let ipc = Ipc::with_event_loop(&*args.flag_ipc, &event_loop.handle())
        .expect("should be able to connect to local unix socket");

    let address: Address = args.flag_contract.parse().expect(
        "should be able to parse address",
    );

    let filter_builder = FilterBuilder::default()
        //.from_block(BlockNumber::Number(0))
        //.to_block(BlockNumber::Latest)
        // .limit(1)
        .address(vec![address]);

    println!("creating transport");
    let transport = api::Eth::new(&ipc);

    // searches over (last_block, block_number-delay]

    loop {
        thread::sleep(time::Duration::from_millis(1000));

        let block_number_fut = transport.block_number();
        let block_number = event_loop.run(block_number_fut).unwrap().low_u64();

        if block_number - delay <= last_block {
            continue;
        }

        println!("New block detected: {}", block_number);

        let filter = filter_builder
            .clone()
            .from_block(BlockNumber::Number(last_block+1))
            .to_block(BlockNumber::Number(block_number-delay))
            .build();

        let logs_fut = transport.logs(&filter);
        let logs = event_loop.run(logs_fut).unwrap();

        for log in logs {
            let block = log.block_number;
            println!("got log {:?}", block);
            let log = new_witness(gen_rawlog(log), &peggy::Peggy::default());
            let keypair = Keypair::from_bytes("test".as_bytes()).expect("frombytes keypair");
            let seq = 0;

            let tx = sign_and_wrap(log, keypair, seq);
        }

        last_block = block_number-delay;
    }
}
