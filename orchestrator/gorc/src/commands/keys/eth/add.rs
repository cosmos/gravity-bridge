use super::show::ShowEthKeyCmd;
use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use bip32;
use k256::pkcs8::ToPrivateKey;
use rand_core::OsRng;
use signatory::FsKeyStore;
use std::path::Path;

#[derive(Command, Debug, Default, Options)]
pub struct AddEthKeyCmd {
    #[options(free, help = "add [name]")]
    pub args: Vec<String>,

    #[options(help = "overwrite existing key")]
    pub overwrite: bool,
}

// Entry point for `gorc keys eth add [name]`
// - [name] required; key name
impl Runnable for AddEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = Path::new(&config.keystore);
        let keystore = FsKeyStore::create_or_open(keystore).expect("Could not open keystore");

        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");
        if let Ok(_info) = keystore.info(&name) {
            if !self.overwrite {
                println!("Key already exists, exiting.");
                return;
            }
        }

        let mnemonic = bip32::Mnemonic::random(&mut OsRng, Default::default());
        println! {"**Important** record this mnemonic in a safe place:"}
        println! {"{}", mnemonic.phrase()};

        let seed = mnemonic.to_seed("");

        let path = config.ethereum.key_derivation_path.trim();
        let path = path
            .parse::<bip32::DerivationPath>()
            .expect("Could not parse derivation path");

        let key = bip32::XPrv::derive_from_path(seed, &path).unwrap();
        let key = k256::SecretKey::from(key.private_key());
        let key = key.to_pkcs8_der().unwrap();
        keystore.store(&name, &key).expect("Could not store key");

        let args = vec![name.to_string()];
        let show_cmd = ShowEthKeyCmd { args };
        show_cmd.run();
    }
}
