use clarity::Address as EthAddress;
use clarity::Signature as EthSignature;
use contact::types::parse_val;
use deep_space::address::Address;
use num256::Uint256;

/// the response we get when querying for a valset confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct ValsetConfirmResponse {
    #[serde(deserialize_with = "parse_val")]
    pub validator: Address,
    #[serde(deserialize_with = "parse_val")]
    pub nonce: Uint256,
    #[serde(deserialize_with = "parse_val", rename = "signature")]
    pub eth_signature: EthSignature,
}

/// a list of validators, powers, and eth addresses at a given block height
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct Valset {
    pub nonce: u64,
    pub powers: Vec<u64>,
    pub eth_addresses: Vec<Option<EthAddress>>,
}

/// a list of validators, powers, and eth addresses at a given block height
/// this version is used by the endpoint to get the data and is then processed
/// by "convert" into ValsetResponse. Making this struct purely internal
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct ValsetUnparsed {
    #[serde(deserialize_with = "parse_val")]
    nonce: u64,
    powers: Vec<String>,
    eth_addresses: Vec<String>,
}

impl ValsetUnparsed {
    pub fn convert(self) -> Valset {
        let mut out = Vec::new();
        let mut powers = Vec::new();
        for maybe_addr in self.eth_addresses.iter() {
            if maybe_addr.is_empty() {
                out.push(None);
            } else {
                match maybe_addr.parse() {
                    Ok(val) => out.push(Some(val)),
                    Err(_e) => out.push(None),
                }
            }
        }
        for power in self.powers.iter() {
            match power.parse() {
                Ok(val) => powers.push(val),
                Err(_e) => powers.push(0),
            }
        }
        Valset {
            nonce: self.nonce,
            powers,
            eth_addresses: out,
        }
    }
}

/// the query struct required to get the valset request sent by a specific
/// validator. This is required because the url encoded get methods don't
/// parse addresses well. So there's no way to get an individual validators
/// address without sending over a json body
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct QueryValsetConfirm {
    pub nonce: String,
    pub address: String,
}
