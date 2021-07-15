use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use std::path;

#[derive(Command, Debug, Default, Options)]
pub struct DeleteEthKeyCmd {
    #[options(free, help = "delete [name]")]
    pub args: Vec<String>,
}

// Entry point for `gorc keys eth delete [name]`
// - [name] required; key name
impl Runnable for DeleteEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);
        let keystore = signatory::FsKeyStore::create_or_open(keystore).unwrap();

        let name = self.args.get(0).expect("name is required");
        let name = name.parse().expect("Could not parse name");
        keystore.delete(&name).unwrap();
    }
}
