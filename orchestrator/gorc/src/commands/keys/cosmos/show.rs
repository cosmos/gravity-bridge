use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use deep_space;
use signatory::FsKeyStore;
use std::path::Path;

#[derive(Command, Debug, Default, Options)]
pub struct ShowCosmosKeyCmd {
    #[options(free, help = "delete [name]")]
    pub args: Vec<String>,
}

// Entry point for `gorc keys cosmos show [name]`
impl Runnable for ShowCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = Path::new(&config.keystore);
        let keystore = FsKeyStore::create_or_open(keystore).unwrap();
        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");

        let key = keystore.load(&name).expect("Could not load key");
        let key = key
            .to_pem()
            .parse::<k256::elliptic_curve::SecretKey<k256::Secp256k1>>()
            .expect("Could not parse key");

        let key = deep_space::utils::bytes_to_hex_str(&key.to_bytes());
        let key = key
            .parse::<deep_space::private_key::PrivateKey>()
            .expect("Could not parse private key");

        let address = key
            .to_address("cosmos")
            .expect("Could not generate public key");

        println!("{}\t{}", name, address)
    }
}