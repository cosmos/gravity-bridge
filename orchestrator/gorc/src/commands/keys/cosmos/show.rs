use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use signatory::FsKeyStore;
use std::path;

#[derive(Command, Debug, Default, Options)]
pub struct ShowCosmosKeyCmd {
    #[options(short = "n", long = "name", help = "show [name]")]
    pub name: String,
}

/// The `gorc keys cosmos show [name]` subcommand: show keys
impl Runnable for ShowCosmosKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to list keys
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();
        let key_name = &self.name.parse().unwrap();
        let show_key = FsKeyStore::info(&keystore, &key_name).unwrap();
        println!("{:?}", show_key)
    }
}
