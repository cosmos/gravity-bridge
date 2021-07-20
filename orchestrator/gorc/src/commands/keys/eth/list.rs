use super::show::ShowEthKeyCmd;
use crate::application::APP;
use abscissa_core::{Application, Command, Options, Runnable};
use std::path;

#[derive(Command, Debug, Default, Options)]
pub struct ListEthKeyCmd {
    #[options(help = "show private key")]
    pub show_private_key: bool,
}

// Entry point for `gorc keys eth list`
impl Runnable for ListEthKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let keystore = path::Path::new(&config.keystore);

        for entry in keystore.read_dir().expect("Could not read keystore") {
            let path = entry.unwrap().path();
            if path.is_file() {
                if let Some(extension) = path.extension() {
                    if extension == "pem" {
                        let name = path.file_stem().unwrap();
                        let name = name.to_str().unwrap();
                        let show_cmd = ShowEthKeyCmd {
                            args: vec![name.to_string()],
                            show_private_key: self.show_private_key,
                            show_name: true,
                        };
                        show_cmd.run();
                    }
                }
            }
        }
    }
}
