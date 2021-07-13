use abscissa_core::{Command, Options, Runnable};
use bip32::{Mnemonic, XPrv};
use rand_core::OsRng;
use signatory::FsKeyStore;
use std::path::Path;
use k256::pkcs8::ToPrivateKey;

#[derive(Command, Debug, Default, Options)]
pub struct AddCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "add private key")]
    pub name: String,
}

/// The `gorc keys cosmos add [name] ` subcommand: add private key & save to disk
impl Runnable for AddCosmosKeyCmd {
    fn run(&self) {
        let mnemonic = Mnemonic::random(&mut OsRng, Default::default());
        println!{"**Important** write this mnemonic in a safe place.\n"}

        println!{"{}", mnemonic.phrase()};
        let seed = mnemonic.to_seed("TREZOR"); // todo: password argument
        let xprv = XPrv::new(&seed).unwrap();
        let private_key_der = k256::SecretKey::from(xprv.private_key()).to_pkcs8_der();
        let private_key_der = private_key_der.unwrap();

        // todo: where the keys go? load from config? for now use /tmp for testing
        let keystore_path = Path::new("/tmp/keystore");
        let keystore = FsKeyStore::create_or_open(keystore_path).unwrap();
        let key_name = &self.name.parse().unwrap();
        keystore.store(key_name, &private_key_der).unwrap();
    }
}
