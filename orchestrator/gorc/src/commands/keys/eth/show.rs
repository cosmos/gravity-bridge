use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct ShowEthKeyCmd {
    #[options(free, help = "show [name]")]
    pub args: Vec<String>,
}

// Entry point for `gorc keys eth show [name]`
impl Runnable for ShowEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let name = self.args.get(0).expect("name is required");
        let key = config.load_clarity_key(name.clone());

        let pub_key = key.to_public_key().expect("Could not build public key");

        println!("{}\t{}", name, pub_key);
    }
}
