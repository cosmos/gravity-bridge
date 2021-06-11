//! contains configuration structs that need to be accessed across crates.

/// Global configuration struct for Gravity bridge tools
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Default)]
pub struct GravityBridgeToolsConfig {
    pub relayer: RelayerConfig,
    pub orchestrator: OrchestratorConfig,
}

/// Relayer configuration options
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq, Default)]
pub struct RelayerConfig {}

/// Orchestrator configuration options
#[derive(Serialize, Deserialize, Debug, PartialEq, Eq)]
pub struct OrchestratorConfig {
    /// If this Orchestrator should run an integrated relayer or not
    pub relayer_enabled: bool,
}

impl Default for OrchestratorConfig {
    fn default() -> Self {
        OrchestratorConfig {
            relayer_enabled: true,
        }
    }
}
