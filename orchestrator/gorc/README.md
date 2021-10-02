# Gorc

Gorc is an application.

## Getting Started

This application is authored using [Abscissa], a Rust application framework.

The `Gorc` application is still under development and the help message isn't very helpful. This is a comprehensive documentation on how to use `Gorc`.

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
cosmos-to-eth [gravity_denom] [amount] [cosmos_key] [eth_dest] [times]
```
The `cosmos-to-eth` command takes the following arguement/flags;

- gravity_denom: The gravity denom.
- amount: Amount to be sent.
- cosmos key: Cosmos private key.
- eth_dest: Ethereum destination address.
- times: 

**deploy:** To deploy an `erc20` contract, run the command below:

```
deploy erc20 [denom] -e [ETHEREUM-KEY]
```
The `deploy` command takes the following arguement/flags;

- denom:
- e:

**eth-to-cosmos:** To send Ethereum to Cosmos, run the command below:

```
eth-to-cosmos [erc20_address] [ethereum_key] [contract_address] [cosmos_dest] [amount] [times]
```

- erc20_address:
- ethereum_key:
- contract_address:
- cosmos_dest:
- amount:
- times:

**keys:**

**orchestrator**

**print-config**

**query**

**sign-delegate-keys**

**tests**

**tx**

For more information, see:

[Documentation]

[Abscissa]: https://github.com/iqlusioninc/abscissa
[Documentation]: https://docs.rs/abscissa_core/
