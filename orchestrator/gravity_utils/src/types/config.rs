//! contains configuration structs that need to be accessed across crates.

/// Global configuration struct for Gravity bridge tools
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Default)]
pub struct GravityBridgeToolsConfig {
    #[serde(default = "RelayerConfig::default")]
    pub relayer: RelayerConfig,
    #[serde(default = "OrchestratorConfig::default")]
    pub orchestrator: OrchestratorConfig,
}

/// Relayer configuration options
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Default)]
pub struct RelayerConfig {}

/// Orchestrator configuration options
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq)]
pub struct OrchestratorConfig {
    /// If this Orchestrator should run an integrated relayer or not
    #[serde(default = "default_relayer_enabled")]
    pub relayer_enabled: bool,
}

fn default_relayer_enabled() -> bool {
    true
}

impl Default for OrchestratorConfig {
    fn default() -> Self {
        OrchestratorConfig {
            relayer_enabled: default_relayer_enabled(),
        }
    }
}
