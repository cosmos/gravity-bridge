use crate::application::APP;
use abscissa_core::{Application, Command, Clap, Runnable};

/// Show an Eth Key
#[derive(Command, Debug, Default, Clap)]
pub struct ShowEthKeyCmd {
    pub args: Vec<String>,

    #[clap(short, long)]
    pub show_private_key: bool,

    #[clap(short = 'n', long)]
    pub show_name: bool,
}

// Entry point for `gorc keys eth show [name]`
impl Runnable for ShowEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let name = self.args.get(0).expect("name is required");
        let key = config.load_clarity_key(name.clone());

        let pub_key = key.to_public_key().expect("Could not build public key");

        if self.show_name {
            print!("{}\t", name);
        }

        if self.show_private_key {
            println!("{}\t{}", pub_key, key);
        } else {
            println!("{}", pub_key);
        }
    }
}
