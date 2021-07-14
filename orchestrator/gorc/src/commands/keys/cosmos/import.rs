use abscissa_core::{Application, Command, Options, Runnable};
use bip32;
use crate::application::APP;
use k256::pkcs8::ToPrivateKey;
use signatory;
use std::path;

#[derive(Command, Debug, Default, Options)]
pub struct ImportCosmosKeyCmd {
    #[options(free, help = "import [name] (mnemonic) (password)")]
    pub args: Vec<String>,

    #[options(help = "overwrite existing key")]
    pub overwrite: bool,
}

// `gorc keys cosmos import [name] (mnemonic) (password)`
// - [name] required; key name
// - (mnemonic) optional; when absent the user will be prompted to enter it
// - (password) optional; when absent the user will be prompted to enter it
impl Runnable for ImportCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();

        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");
        if let Ok(_info) = keystore.info(&name) {
            if !self.overwrite {
                println!("Key already exists, exiting.");
                return;
            }
        }

        let mnemonic = match self.args.get(1) {
            Some(mnemonic) => mnemonic.clone(),
            None => rpassword::read_password_from_tty(Some("Mnemonic: ")).unwrap(),
        };

        let password = match self.args.get(2) {
            Some(password) => password.clone(),
            None => rpassword::read_password_from_tty(Some("Password: ")).unwrap(),
        };

        let mnemonic = bip32::Mnemonic::new(mnemonic.trim_end(), Default::default()).unwrap();

        let seed = mnemonic.to_seed(&password);
        let path = "m/44'/118'/0'/0/0".parse::<bip32::DerivationPath>().unwrap();
        let key = bip32::XPrv::derive_from_path(seed, &path).unwrap();
        let key = k256::SecretKey::from(key.private_key());
        let key = key.to_pkcs8_der().unwrap();

        keystore.store(&name, &key).expect("Could not store key");
    }
}