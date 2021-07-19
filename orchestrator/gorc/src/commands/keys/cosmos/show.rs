use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct ShowCosmosKeyCmd {
    #[options(free, help = "show [name]")]
    pub args: Vec<String>,
}

// Entry point for `gorc keys cosmos show [name]`
impl Runnable for ShowCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let name = self.args.get(0).expect("name is required");
        let key = config.load_deep_space_key(name.clone());

        let address = key
            .to_address(config.cosmos.prefix.trim())
            .expect("Could not generate public key");

        println!("{}\t{}", name, address)
    }
}