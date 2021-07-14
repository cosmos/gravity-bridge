use abscissa_core::{Command, Options, Runnable};
use signatory::FsKeyStore;
use std::path::Path;

#[derive(Command, Debug, Default, Options)]
pub struct ShowCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "show [name]")]
    pub name: String,
}

/// The `gorc keys cosmos show [name]` subcommand: show keys
impl Runnable for ShowCosmosKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to list keys
        let keystore_path = Path::new("/tmp/keystore");
        let keystore = FsKeyStore::create_or_open(keystore_path).unwrap();
        let key_name = &self.name.parse().unwrap();
        let key_info = keystore.info(&key_name).unwrap();
        let show_key = FsKeyStore::info(&keystore, &key_name).unwrap();
        println!("{:?}", show_key)
    }
}
