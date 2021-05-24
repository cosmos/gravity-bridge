//! Handles configuration structs + saving and loading for Gravity bridge tools

use std::{
    fs::{self, create_dir},
    path::PathBuf,
    process::exit,
};

use crate::args::InitOpts;

/// The name of the config file
pub const CONFIG_NAME: &str = "config.toml";

/// Global configuration struct for Gravity bridge tools
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Default)]
pub struct GravityBridgeToolsConfig {}

/// Relayer configuration options
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Default)]
pub struct Relayer {}

/// Orchestrator configuration options
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Default)]
pub struct Orchestrator {}

/// Creates the config directory and default config file if it does
/// not already exist
pub fn init_config(_init_ops: InitOpts, home_dir: PathBuf) {
    if home_dir.exists() {
        warn!(
            "The Gravity bridge tools config folder {} already exists!",
            home_dir.to_str().unwrap()
        );
        warn!("You can delete this folder and run init again, you will lose any keys or other config data!");
    } else {
        create_dir(home_dir.clone()).expect("Failed to create config directory!");

        fs::write(home_dir.with_file_name(CONFIG_NAME), get_default_config())
            .expect("Unable to write config file");
    }
}

/// Loads the default config from the default-config.toml file
/// done at compile time and is included in the binary
/// This is done so that we can have hand edited and annotated
/// config
fn get_default_config() -> String {
    include_str!("default-config.toml").to_string()
}

/// Load the config file, this operates at runtime
pub fn load_config(home_dir: PathBuf) -> GravityBridgeToolsConfig {
    let config = fs::read_to_string(home_dir.with_file_name(CONFIG_NAME))
        .expect("Could not find config file! Run `gbt init`");
    match toml::from_str(&config) {
        Ok(v) => v,
        Err(e) => {
            error!("Invalid config! {:?}", e);
            exit(1);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    /// Test that the config is both valid toml for the struct and that it's values are
    /// equal to the default values of the config.
    #[test]
    fn test_default_config() {
        let res: GravityBridgeToolsConfig = toml::from_str(&get_default_config()).unwrap();
        assert_eq!(res, GravityBridgeToolsConfig::default())
    }
}
