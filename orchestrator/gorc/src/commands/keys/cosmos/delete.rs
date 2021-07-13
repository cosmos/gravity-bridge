use abscissa_core::{Command, Options, Runnable};
use signatory::FsKeyStore;
use std::path::Path;

#[derive(Command, Debug, Default, Options)]
pub struct DeleteCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "delete key [name]")]
    pub name: String,
}

/// The `gork keys cosmos delete [name] ` subcommand: delete the given key
impl Runnable for DeleteCosmosKeyCmd {
    fn run(&self) {
        let keystore_path = Path::new("/tmp/keystore");
        let keystore = FsKeyStore::create_or_open(keystore_path).unwrap();
        pub const EXAMPLE_KEY: &str = "example-key";
        let key_name = EXAMPLE_KEY.parse().unwrap();
        let _delete_key = FsKeyStore::delete(&keystore, &key_name).unwrap();
    }
}
