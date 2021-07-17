use abscissa_core::{Application, Command, Options, Runnable};
use crate::application::APP;
use std::path;

#[derive(Command, Debug, Default, Options)]
pub struct RenameCosmosKeyCmd {
    #[options(free, help = "rename [name] [new_name]")]
    pub args: Vec<String>,

    #[options(help = "overwrite existing key")]
    pub overwrite: bool,
}

/// The `gorc keys cosmos rename [name] [new-name]` subcommand: show keys
impl Runnable for RenameCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();

        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");

        let new_name = self.args.get(1).expect("new_name is required");
        let new_name = new_name.parse().expect("Could not parse new_name");
        if let Ok(_info) = keystore.info(&new_name) {
            if !self.overwrite {
                println!("Key already exists, exiting.");
                return;
            }
        }

        let key = keystore.load(&name).expect("Could not load key");
        keystore.store(&new_name, &key).unwrap();
        keystore.delete(&name).unwrap();
    }
}
