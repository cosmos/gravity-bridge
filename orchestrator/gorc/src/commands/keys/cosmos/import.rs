use abscissa_core::{Command, Options, Runnable};
use bip32::{Mnemonic, XPrv};

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
        dbg!{mnemonic.phrase()};

        let seed = mnemonic.to_seed("TREZOR"); // TODO: password argument
        let root_key = XPrv::new(&seed).unwrap();
        let expected_key: XPrv = "xprv9s21ZrQH143K3Y1sd2XVu9wtqxJRvybCfAetjUrMMco6r3v9qZTBeXiBZkS8JxWbcGJZyio8TrZtm6pkbzG8SYt1sxwNLh3Wx7to5pgiVFU".parse().unwrap();
        assert_eq!(root_key.secret_key().to_bytes(), expected_key.secret_key().to_bytes());
        println!("OK!")
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