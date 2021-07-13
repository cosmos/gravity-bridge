### Conceptual Goals:
* Have an easily understood network layout based on docker or other infrastructure tooling
* Have a one click deploy of a full network
* Allow definition of deployment of a network in certain states or based on defined preconditions

### Deliverable Goals:
* Docker images for each component that are self-contained and defined to run as production assets
* Test suite definition tool that uses said image in combination with arguments/files for the scenario

#### Current state:
* Orchestrator
  * Defined docker image with rust binaries.
  * No files required, all settings brought in via image environment variables
* Gravity module
  * Docker build requires keys signed and rotated before build
  * Environment variables/arguments are unused, all configuration comes from files
* Ethereum
  * Most likely to be changed or need specific starting states, least important to be deployable image because it isn't our product
  * Can be implement with geth+genesis file, or a hardhat image