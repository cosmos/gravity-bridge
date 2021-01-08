//! This crate contains various components and utilities for interacting with the Peggy Cosmos module. Primarily
//! Extensions to Althea's 'deep_space' Cosmos transaction library to allow it to send Peggy module specific messages
//! parse Peggy module specific endpoints and generally interact with the multitude of Peggy specific functionality
//! that's part of the Cosmos module.

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate log;

pub mod messages;
pub mod query;
pub mod send;
pub mod utils;
