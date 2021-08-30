use serde::{Deserialize, Serialize};
use signatory::FsKeyStore;
use std::net::SocketAddr;
use std::path::Path;

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(default, deny_unknown_fields)]
pub struct GorcConfig {
    pub keystore: String,
    pub gravity: GravitySection,
    pub ethereum: EthereumSection,
    pub cosmos: CosmosSection,
    pub metrics: MetricsSection,
}

impl GorcConfig {
    fn load_secret_key(&self, name: String) -> k256::elliptic_curve::SecretKey<k256::Secp256k1> {
        let keystore = Path::new(&self.keystore);
        let keystore = FsKeyStore::create_or_open(keystore).expect("Could not open keystore");
        let name = name.parse().expect("Could not parse name");
        let key = keystore.load(&name).expect("Could not load key");
        return key.to_pem().parse().expect("Could not parse pem");
    }

    pub fn load_clarity_key(&self, name: String) -> clarity::PrivateKey {
        let key = self.load_secret_key(name).to_bytes();
        return clarity::PrivateKey::from_slice(&key).expect("Could not convert key");
    }

    pub fn load_deep_space_key(&self, name: String) -> deep_space::private_key::PrivateKey {
        let key = self.load_secret_key(name).to_bytes();
        let key = deep_space::utils::bytes_to_hex_str(&key);
        return key.parse().expect("Could not parse private key");
    }
}

impl Default for GorcConfig {
    fn default() -> Self {
        Self {
            keystore: "/tmp/keystore".to_owned(),
            gravity: GravitySection::default(),
            ethereum: EthereumSection::default(),
            cosmos: CosmosSection::default(),
            metrics: MetricsSection::default(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(default, deny_unknown_fields)]
pub struct GravitySection {
    pub contract: String,
    pub fees_denom: String,
}

impl Default for GravitySection {
    fn default() -> Self {
        Self {
            contract: "0x0000000000000000000000000000000000000000".to_owned(),
            fees_denom: "stake".to_owned(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(default, deny_unknown_fields)]
pub struct EthereumSection {
    pub key_derivation_path: String,
    pub rpc: String,
}

impl Default for EthereumSection {
    fn default() -> Self {
        Self {
            key_derivation_path: "m/44'/60'/0'/0/0".to_owned(),
            rpc: "http://localhost:8545".to_owned(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(default, deny_unknown_fields)]
pub struct CosmosSection {
    pub key_derivation_path: String,
    pub grpc: String,
    pub prefix: String,
    pub gas_price: GasPrice,
}

impl Default for CosmosSection {
    fn default() -> Self {
        Self {
            key_derivation_path: "m/44'/118'/0'/0/0".to_owned(),
            grpc: "http://localhost:9090".to_owned(),
            prefix: "cosmos".to_owned(),
            gas_price: GasPrice::default(),
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(default, deny_unknown_fields)]
pub struct GasPrice {
    pub amount: f64,
    pub denom: String,
}

impl Default for GasPrice {
    fn default() -> Self {
        Self {
            amount: 0.001,
            denom: "stake".to_owned(),
        }
    }
}

impl GasPrice {
    pub fn as_tuple(&self) -> (f64, String) {
        (self.amount, self.denom.to_owned())
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(default, deny_unknown_fields)]
pub struct MetricsSection {
    pub listen_addr: SocketAddr,
}

impl Default for MetricsSection {
    fn default() -> Self {
        Self {
            listen_addr: "127.0.0.1:3000".parse().unwrap(),
        }
    }
}
