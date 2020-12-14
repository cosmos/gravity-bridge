//! Test runner is a testing script for the Peggy Cosmos module. It is built in Rust rather than python or bash
//! to maximize code and tooling shared with the validator-daemon and relayer binaries.
//! For the core of the logic please see the test_cases module

// there are several binaries for this crate if we allow dead code on all of them
// we will see functions not used in one binary as dead code. In order to fix that
// we forbid dead code in all but the 'main' binary
#![allow(dead_code)]

#[macro_use]
extern crate log;
#[macro_use]
extern crate lazy_static;

mod batch_relaying;
mod ethereum_event_watcher;
mod main_loop;
mod test_cases;
mod valset_relaying;

use test_cases::run_peggy_test_cases;

#[actix_rt::main]
async fn main() {
    env_logger::init();
    info!("Staring Peggy test-runner");
    run_peggy_test_cases().await;
}
