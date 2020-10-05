//! This module is used to interact with the Peggy module on the Cosmos chain
//! this includes transactions and rest endpoints.

pub mod msgs;
pub mod query;
pub mod send;
#[cfg(test)]
pub mod tests;
pub mod types;
