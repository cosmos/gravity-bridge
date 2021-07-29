use abscissa_core::{Runnable, Command,Options};

#[derive(Command, Debug, Default, Options)]
pub struct PrintConfigCmd {}

impl Runnable for PrintConfigCmd {
    fn run(&self) {
        let config = crate::config::GorcConfig::default();
        print!("{}", toml::to_string(&config).unwrap());
    }
}