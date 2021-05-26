use abscissa_core::{Command, Options, Runnable};
use bip32::{Mnemonic, XPrv};
use pkcs8::ToPrivateKey;
use rand_core::OsRng;
use signatory::keystore::FsKeyStore;
use std::path::Path;

#[derive(Command, Debug, Default, Options)]
pub struct AddCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "add private key [name]")]
    pub name: String,
}

/// The `gork keys cosmos add [name] ` subcommand: add private key & save to disk
impl Runnable for AddCosmosKeyCmd {
    fn run(&self) {
        let mnemonic = Mnemonic::random(&mut OsRng, Default::default());
        println!{"**Important** write this mnemonic in a safe place.\n"}

        println!{"{}", mnemonic.phrase()};
        let seed = mnemonic.to_seed("TREZOR"); // todo: password argument
        let xprv = XPrv::new(&seed).unwrap();
        let private_key_der = k256::SecretKey::from(xprv.private_key()).to_pkcs8_der();

        // todo: where the keys go? load from config? for now use /tmp for testing
        let keystore_path = Path::new("/tmp/keystore");
        if !keystore_path.exists() {
            FsKeyStore::create(keystore_path).unwrap();
        }
        let keystore = FsKeyStore::open(keystore_path).unwrap();
        keystore.store(&self.name, &private_key_der).unwrap();
    }
}
