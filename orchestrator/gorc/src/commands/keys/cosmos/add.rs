use abscissa_core::{Application, Command, Options, Runnable};
use bip32;
use crate::application::APP;
use k256::pkcs8::ToPrivateKey;
use rand_core::OsRng;
use std::path;

#[derive(Command, Debug, Default, Options)]
pub struct AddCosmosKeyCmd {
    #[options(free, help = "add [name] (password)")]
    pub args: Vec<String>,

    #[options(help = "overwrite existing key")]
    pub overwrite: bool,
}

// `gorc keys cosmos add [name] (password)`
// - [name] required; key name
// - (password) optional; when absent the user will be prompted to enter it
impl Runnable for AddCosmosKeyCmd {
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

        let password = match self.args.get(1) {
            Some(password) => password.clone(),
            None => rpassword::read_password_from_tty(Some("Password: ")).unwrap(),
        };

        let mnemonic = bip32::Mnemonic::random(&mut OsRng, Default::default());
        println! {"**Important** record this mnemonic in a safe place:"}
        println! {"{}", mnemonic.phrase()};

        let key = bip32::XPrv::new(mnemonic.to_seed(&password)).unwrap();
        let key = k256::SecretKey::from(key.private_key());
        let key = key.to_pkcs8_der().unwrap();

        keystore.store(&name, &key).unwrap();
    }
}