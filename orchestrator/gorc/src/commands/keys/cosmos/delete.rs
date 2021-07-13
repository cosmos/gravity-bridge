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
        // Path where key is stored.
        let keystore_path = Path::new("/tmp/keystore");
        let keystore = FsKeyStore::create_or_open(keystore_path).unwrap();
        // Collect key name from args.
        let key_name = &self.name.parse().unwrap();
        // Delete keyname after locating file from path and key name.
        let _delete_key = FsKeyStore::delete(&keystore, &key_name).unwrap();
    }
}
