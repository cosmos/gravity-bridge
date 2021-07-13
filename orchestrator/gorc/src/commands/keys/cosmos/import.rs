use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use bip32;
use k256::pkcs8::ToPrivateKey;
use signatory;
use std::path;

#[derive(Command, Debug, Default, Options)]
pub struct ImportCosmosKeyCmd {
    #[options(free, help = "import [name] (mnemonic) (password)")]
    pub args: Vec<String>,
}

// Args:
// - name is required
// - mnemonic is optional; when absent the user will be prompted to enter it
// - password is optional; when absent the user will be prompted to enter it
impl Runnable for ImportCosmosKeyCmd {
    fn run(&self) {
        // TODO(levi) make sure there's at least one arg and no more than 3

        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");

        let mnemonic = match self.args.get(1) {
            Some(mnemonic) => mnemonic.clone(),
            None => rpassword::read_password_from_tty(Some("Mnemonic: ")).unwrap(),
        };

        let password = match self.args.get(2) {
            Some(password) => password.clone(),
            None => rpassword::read_password_from_tty(Some("Password: ")).unwrap(),
        };

        let config = APP.config();

        let mnemonic = bip32::Mnemonic::new(mnemonic.trim_end(), Default::default()).unwrap();
        let key = bip32::XPrv::new(mnemonic.to_seed(&password)).unwrap();
        let key = k256::SecretKey::from(key.private_key());
        let key = key.to_pkcs8_der().unwrap();

        let keystore = path::Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();
        keystore.store(&name, &key).expect("Could not store key");
    }
}