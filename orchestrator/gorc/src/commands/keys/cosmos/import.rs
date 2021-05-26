use abscissa_core::{Command, Options, Runnable};
use bip32::{Mnemonic, XPrv};
use pkcs8::ToPrivateKey;
use signatory::keystore::FsKeyStore;
use std::path::Path;

#[derive(Command, Debug, Default, Options)]
pub struct ImportCosmosKeyCmd {
    #[options(
        short = "n",
        long = "name",
        help = "import private key [name] [mnemnoic]"
    )]
    pub name: String,

    #[options(
        short = "m",
        long = "mnemnoic",
        help = "import private key [name] [mnemnoic]"
    )]
    pub mnemnoic: String,
}

/// The `gork keys cosmos import [name] [mnemnoic]` subcommand: import key
impl Runnable for ImportCosmosKeyCmd {
    fn run(&self) {
        let phrase = rpassword::read_password_from_tty(Some("Mnemonic: ")).unwrap();
        let mnemonic = Mnemonic::new(phrase.trim_end(), Default::default()).unwrap();
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

// #[cfg(test)]
// mod tests {
//     use bip32::{Mnemonic, Language, XPrv};
//
//     fn test_vector() {
//         let password = "TREZOR";
//
//
//             let mnemonic = mnemonic::Phrase::new(vector.phrase, Language::default()).unwrap();
//             assert_eq!(&vector.seed, mnemonic.to_seed(password).as_bytes());
//         }
//     }
// }