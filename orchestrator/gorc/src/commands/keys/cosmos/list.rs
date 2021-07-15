use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use std::path::Path;

#[derive(Command, Debug, Default, Options)]
pub struct ListCosmosKeyCmd {}

/// The `gorc keys cosmos list` subcommand: list keys
impl Runnable for ListCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = Path::new(&config.keystore);
        let new_keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();

        for entry in keystore.read_dir().expect("Could not read keystore") {
            let path = entry.unwrap().path();
            if path.is_file() {
                if let Some(extension) = path.extension() {
                    if extension == "pem" {
                        let name = path.file_stem().unwrap();
                        let name = name.to_str().unwrap();
                        let name = name.parse().expect("Could not parse name");

                        let key = new_keystore.load(&name).expect("Could not load key");
                        let key = key
                            .to_pem()
                            .parse::<k256::elliptic_curve::SecretKey<k256::Secp256k1>>()
                            .expect("Could not parse key");
                        let key = clarity::PrivateKey::from_slice(&key.to_bytes()).unwrap();
                        let pub_key = key.to_public_key().unwrap();
                        println!("{}\t{}", name, pub_key)
                    }
                }
            }
        }
    }
}
