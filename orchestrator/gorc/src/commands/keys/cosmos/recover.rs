use super::show::ShowCosmosKeyCmd;
use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use k256::pkcs8::ToPrivateKey;
use signatory::FsKeyStore;
use std::path;

#[derive(Command, Debug, Default, Options)]
pub struct RecoverCosmosKeyCmd {
    #[options(free, help = "recover [name] (bip39-mnemonic)")]
    pub args: Vec<String>,

    #[options(help = "overwrite existing key")]
    pub overwrite: bool,
}

// `gorc keys cosmos recover [name] (bip39-mnemonic)`
// - [name] required; key name
// - (bip39-mnemonic) optional; when absent the user will be prompted to enter it
impl Runnable for RecoverCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore =
            FsKeyStore::create_or_open(keystore).expect("Could not open keystore");

        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");
        if let Ok(_info) = keystore.info(&name) {
            if !self.overwrite {
                eprintln!("Key already exists, exiting.");
                return;
            }
        }

        let mnemonic = match self.args.get(1) {
            Some(mnemonic) => mnemonic.clone(),
            None => rpassword::read_password_from_tty(Some(
                "> Enter your bip39-mnemonic:\n",
            ))
            .expect("Could not read mnemonic"),
        };

        let mnemonic = bip32::Mnemonic::new(mnemonic.trim(), Default::default())
            .expect("Could not parse mnemonic");

        let seed = mnemonic.to_seed("");

        let path = config.cosmos.key_derivation_path.clone();
        let path = path
            .parse::<bip32::DerivationPath>()
            .expect("Could not parse derivation path");

        let key = bip32::XPrv::derive_from_path(seed, &path).expect("Could not derive key");
        let key = k256::SecretKey::from(key.private_key());
        let key = key
            .to_pkcs8_der()
            .expect("Could not PKCS8 encod private key");

        keystore.store(&name, &key).expect("Could not store key");

        let args = vec![name.to_string()];
        let show_cmd = ShowCosmosKeyCmd { args };
        show_cmd.run();
    }
}
