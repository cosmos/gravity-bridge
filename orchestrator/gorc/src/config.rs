use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct GorcConfig {
    pub gravity: GravitySection,
    pub ethereum: EthereumSection,
    pub cosmos: CosmosSection,
}

impl Default for GorcConfig {
    fn default() -> Self {
        Self {
            gravity: GravitySection::default(),
            ethereum: EthereumSection::default(),
            cosmos: CosmosSection::default(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct GravitySection {
    pub contract: String,
}

impl Default for GravitySection {
    fn default() -> Self {
        Self {
            contract: "0x6b175474e89094c44da98b954eedeac495271d0f".to_owned(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct EthereumSection {
    pub key: String,
    pub rpc: String,
}

impl Default for EthereumSection {
    fn default() -> Self {
        Self {
            key: "testkey".to_owned(),
            rpc: "http://localhost:8545".to_owned(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
pub struct CosmosSection {
    pub key: String,
    pub grpc: String,
    pub prefix: String,
}

impl Default for CosmosSection {
    fn default() -> Self {
        Self {
            key: "testkey".to_owned(),
            grpc: "http://localhost:9090".to_owned(),
            prefix: "cosmos".to_owned(),
        }
    }
}
