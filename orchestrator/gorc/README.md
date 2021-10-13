# Gorc

Gorc is an application for the gravity <-> ethereum bridge.

## Getting Started

This application is authored using [Abscissa], a Rust application framework.

The `Gorc` application is still under development. This is a comprehensive documentation on how to use `Gorc`.

### Gorc subcommands

| Subcommand        | Description                                                 |
| ----------------- | ----------------------------------------------------------- |
| cosmos-to-eth     | This command, send Cosmos to Ethereum                       |
| deploy            | Provides tools for contract deployment                      |
| eth-to-cosmos     | This command, send Ethereum to Cosmos                       |
| help              | Help command to get usage information                       |
| keys              | Key management commands for Ethereum and Cosmos             |
| orchestrator      | Management commannds for the orchestrator                   |
| print-config      | Command for printing configurations                         |
| query             | Command to query state on either ethereum or cosmos chains  |
| sign-delegate-keys| Sign delegate keys command                                  |
| tests             | Command to run tests against configured chains              |
| tx                | Create transactions on either ethereum or cosmos chains     |
| version           | Display version information                                 |

**cosmos-to-eth:** To send Cosmos to Ethereum, run the command below:

```
gorc cosmos-to-eth [gravity_denom] [amount] [cosmos_key] [eth_dest] [times]
```
The `cosmos-to-eth` command takes the following argument/flags;

- gravity_denom: The gravity denom.
- amount: Amount to be sent.
- cosmos key: Cosmos private key.
- eth_dest: Ethereum destination address.
- times: The times.

**deploy:** To deploy an `erc20` contract, run the command below:

```
gorc deploy erc20 [denom] -e [ETHEREUM-KEY]
```
The `deploy` command takes the following argument/flags;

- denom:  The denom
- e: Eth flag

**eth-to-cosmos:** To send Ethereum to Cosmos, run the command below:

```
gorc eth-to-cosmos [erc20_address] [ethereum_key] [contract_address] [cosmos_dest] [amount] [times]
```
The `eth-to-cosmos` command takes the following argument/flags;

- erc20_address: The `erc20` address.
- ethereum_key: The Ethereum private key.
- contract_address: The contract address.
- cosmos_dest: The Cosmos destination.
- amount: Amount to be sent.
- times: The times

**keys:** To manage keys in the Cosmos and Ethereum chain, run any of the commands below:

```
// Create a new cosmos key
gorc keys cosmos add [name]

// Delete a cosmos key
gorc keys cosmos delete [name]

// List all cosmos keys
gorc keys cosmos list

// Recover a cosmos key
gorc keys cosmos recover [name] (bip39-mnemonic)

// List a particular cosmos key
gorc keys cosmos show [name]

// Create a new eth key
gorc keys eth add [name]

// Delete an eth key
gorc keys eth delete [name]

// Import an eth key
gorc keys eth import [name] (private-key)

// List all eth keys
gorc keys eth list

// Recover an eth key
gorc keys eth recover [name] (bip39-mnemonic)

// Rename an eth key
gorc keys eth rename [name] [new-name]

// List a particular eth key
gorc keys eth show [name]
```

The `keys` command takes the following argument/flags;

- name: The key name.
- bip39-mnemonic: Mnemonic seed for generating deterministic keys.
- private-key: Ethereum private key.

**orchestrator:** To start the orchestrator, run the command below:

```
gorc orchestrator start
```

**print-config:** To print the config file in your console, run the command below.

```
gorc print-config
```

**query:** To query either the ethereum or cosmos chain, run any of the commands below:

```
// Query cosmos balance
gorc query cosmos balance [key-name]

// Query cosmos gravity keys
gorc query cosmos gravity-keys [key-name]

// Query eth balance
gorc query eth balance [key-name]

// Query eth contract
gorc query eth contract
```

The `query` command takes the following argument/flags;

- key-name: The key name

**sign-delegate-keys:** To sign delegate keys, run the command below:

```
gorc sign-delegate-key [ethereum-key-name] [validator-address] (nonce)
```

The `sign-delegate-keys` command takes the following argument/flags;

- ethereum-key-name: The Ethereum key name.
- validator-address: The validator address.
- nonce: The nonce.

**tx:** To create transactions on either ethereum or cosmos chains, run any of the commands below:

```
// Send to Ethereum
gorc tx cosmos send-to-eth [from-cosmos-key] [to-eth-addr] [erc20-coin] [[--times=int]]

// Send
gorc tx cosmos send [from-key] [to-addr] [coin-amount]

// Send to Cosmos
gorc tx eth send-to-cosmos [from-eth-key][to-cosmos-addr] [erc20 conract] [erc20 amount] [[--times=int]]

// Send
gorc tx eth send [from-key] [to-addr] [amount] [token-contract]
```

## Note
`[]` means a free argument, `()` means a flag. For instance, `gorc sign-delegate-key [ethereum-key-name] [validator-address] (nonce)` translates to `gorc sign-delegate-key ethereum_key_name validator_address --nonce`.

For more information, see:

[Documentation]

[Abscissa]: https://github.com/iqlusioninc/abscissa
[Documentation]: https://docs.rs/abscissa_core/
