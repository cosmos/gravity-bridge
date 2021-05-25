use abscissa_core::{Command, Options, Runnable};
use bip32::{Mnemonic, ExtendedSecretKey};


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

        let seed = mnemonic.to_seed("");
        let root_key = ExtendedSecretKey::new(seed.as_bytes());
    }
}

