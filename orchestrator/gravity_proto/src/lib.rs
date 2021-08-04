//! This crate provides Gravity proto definitions in Rust and also re-exports cosmos_sdk_proto for use by downstream
//! crates. By default around a dozen proto files are generated and places into the prost folder. We could then proceed
//! to fix up all these files and use them as the required dependencies to the Gravity file, but we chose instead to replace
//! those paths with references ot upstream cosmos-sdk-proto and delete the other files. This reduces cruft in this repo even
//! if it does make for a somewhat more confusing proto generation process.

pub use cosmos_sdk_proto;
pub mod gravity {
    include!("prost/gravity.v1.rs");
}

use bytes::BytesMut;
use prost::Message;
use prost_types::Any;

pub trait ToAny {
    fn to_any(&self) -> Option<prost_types::Any>
    where
        Self: prost::Message;
}

impl ToAny for gravity::BatchExecutedEvent {
    fn to_any(&self) -> Option<prost_types::Any> {
        let mut buf = BytesMut::with_capacity(self.encoded_len());
        self.encode(&mut buf).expect("encoding failed");
        Some(Any {
            type_url: "/gravity.v1.BatchExecutedEvent".into(),
            value: buf.to_vec(),
        })
    }
}

impl ToAny for gravity::BatchTxConfirmation {
    fn to_any(&self) -> Option<prost_types::Any> {
        let mut buf = BytesMut::with_capacity(self.encoded_len());
        self.encode(&mut buf).expect("encoding failed");
        Some(Any {
            type_url: "/gravity.v1.BatchTxConfirmation".into(),
            value: buf.to_vec(),
        })
    }
}

impl ToAny for gravity::ContractCallExecutedEvent {
    fn to_any(&self) -> Option<prost_types::Any> {
        let mut buf = BytesMut::with_capacity(self.encoded_len());
        self.encode(&mut buf).expect("encoding failed");
        Some(Any {
            type_url: "/gravity.v1.ContractCallExecutedEvent".into(),
            value: buf.to_vec(),
        })
    }
}

impl ToAny for gravity::ContractCallTxConfirmation {
    fn to_any(&self) -> Option<prost_types::Any> {
        let mut buf = BytesMut::with_capacity(self.encoded_len());
        self.encode(&mut buf).expect("encoding failed");
        Some(Any {
            type_url: "/gravity.v1.ContractCallTxConfirmation".into(),
            value: buf.to_vec(),
        })
    }
}

impl ToAny for gravity::Erc20DeployedEvent {
    fn to_any(&self) -> Option<prost_types::Any> {
        let mut buf = BytesMut::with_capacity(self.encoded_len());
        self.encode(&mut buf).expect("encoding failed");
        Some(Any {
            type_url: "/gravity.v1.ERC20DeployedEvent".into(),
            value: buf.to_vec(),
        })
    }
}

impl ToAny for gravity::SendToCosmosEvent {
    fn to_any(&self) -> Option<prost_types::Any> {
        let mut buf = BytesMut::with_capacity(self.encoded_len());
        self.encode(&mut buf).expect("encoding failed");
        Some(Any {
            type_url: "/gravity.v1.SendToCosmosEvent".into(),
            value: buf.to_vec(),
        })
    }
}

impl ToAny for gravity::SignerSetTxExecutedEvent {
    fn to_any(&self) -> Option<prost_types::Any> {
        let mut buf = BytesMut::with_capacity(self.encoded_len());
        self.encode(&mut buf).expect("encoding failed");
        Some(Any {
            type_url: "/gravity.v1.SignerSetTxExecutedEvent".into(),
            value: buf.to_vec(),
        })
    }
}

impl ToAny for gravity::SignerSetTxConfirmation {
    fn to_any(&self) -> Option<prost_types::Any> {
        let mut buf = BytesMut::with_capacity(self.encoded_len());
        self.encode(&mut buf).expect("encoding failed");
        Some(Any {
            type_url: "/gravity.v1.SignerSetTxConfirmation".into(),
            value: buf.to_vec(),
        })
    }
}