use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use clarity;
use std::path;

#[derive(Command, Debug, Default, Options)]
pub struct ShowEthKeyCmd {
    #[options(free, help = "show [name]")]
    pub args: Vec<String>,
}

/// The `gorc keys eth show [name]` subcommand: show keys
impl Runnable for ShowEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();

        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");

        let key = keystore.load(&name).expect("Could not load key");
        let key = key
            .to_pem()
            .parse::<k256::elliptic_curve::SecretKey<k256::Secp256k1>>()
            .expect("Could not parse key");
        let key = clarity::PrivateKey::from_slice(&key.to_bytes()).unwrap();

        let pub_key = key.to_public_key().unwrap();
        println!("{}", pub_key);
    }
}
