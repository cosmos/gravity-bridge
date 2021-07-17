use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use signatory::FsKeyStore;
use std::path::Path;

#[derive(Command, Debug, Default, Options)]
pub struct DeleteCosmosKeyCmd {
    #[options(free, help = "delete [name]")]
    pub args: Vec<String>,
}

/// The `gork keys cosmos delete [name] ` subcommand: delete the given key
impl Runnable for DeleteCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        // Path where key is stored.
        let keystore = Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();
        // Collect key name from args.
        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");
        // Delete keyname after locating file from path and key name.
        let _delete_key = FsKeyStore::delete(&keystore, &name).unwrap();
    }
}
