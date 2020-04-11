## Initialization

First, initialize a chain and create accounts:

```bash
# Initialize the genesis.json file that will help you to bootstrap the network
ebd init local --chain-id=peggy

# Configure your cli to use the keyring-backend test so that you don't need to enter a password
# **_NOTE_** Don't use the test backend on production
ebcli config keyring-backend test

# Add some more configuration to avoid using flags
ebcli config chain-id peggy
ebcli config trust-node true
ebcli config indent true
ebcli config output json

# Create a key to hold your validator account and for another test account
ebcli keys add validator
ebcli keys add testuser

# Initialize the genesis account and transaction
ebd add-genesis-account $(ebcli keys show validator -a) 1000000000stake,1000000000atom

# Create genesis transaction
ebd gentx --name validator --keyring-backend test

# Collect genesis transaction
ebd collect-gentxs

# Now its safe to start `ebd`
ebd start
```

## Testing the application

Once you've initialized the application and started the Bridge blockchain with `ebd start`, you can test the available cli commands. They include sending tokens between accounts, querying accounts, claim creation, token burning, and token locking. Once the Relayer is running, you'll be able to submit new burning/locking txs to the chain using these commands.

First, we'll test sending a random token in another terminal window.

```bash
# In another terminal window, send 10 stake tokens from the validator to the testuser
ebcli tx send validator $(ebcli keys show testuser -a) 10stake --yes

# Confirm token balances have changed appropriately
ebcli query account $(ebcli keys show validator -a)
ebcli query account $(ebcli keys show testuser -a)

# Confirm your validator was created correctly, and has become Bonded
ebcli query staking validators

# See more information for any of the ethbridge commands with --help
ebcli tx ethbridge COMMAND --help


# Now you can simulate the process of moving assets between an EVM based chain and the peggy chain. Since there
# is currently no running EVM chain, we will only see the peggy chain update based on fictional claims from
# the validator. Normally these claims are made automatically by the validator when they are operating the
# relayer and real events are witnessed on a running EVM chain.


# Create a bridge lock claim (Ethereum prophecies are stored on the blockchain with an identifier created by
# concatenating the nonce and sender address). Since there is no EVM chain running at the moment, we use the
# address "0x30753E4A8aad7F8597332E813735Def5dD395028" for the bridge-registry-contract because this will be
# the address created the first time you run a local EVM chain using the instructions found in ./setup-eth-local.md.
# For the ethereum-sender-address we will use the vanity address "0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9"
# to make it easier to distinguish. The token-contract is all "0"s because we are transfering the native EVM
# token (ether) rather than an ERC-20.

# See the help for the ethbridge create claim function
ebcli tx ethbridge create-claim --help
# ebcli tx ethbridge create-claim [bridge-registry-contract] [nonce] [symbol] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount] [claim-type] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address] [flags]
ebcli tx ethbridge create-claim 0x30753E4A8aad7F8597332E813735Def5dD395028 0 eth 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 3eth lock --token-contract-address=0x0000000000000000000000000000000000000000 --ethereum-chain-id=3 --from=validator --yes

# You can check the transaction and message were proccessed successfully by querying the transaction hash
# that was just generated using the following command
ebcli q tx TXHASH

# Then read the prophecy to confirm it was created with the claim added
# ebcli query ethbridge prophecy [bridge-registry-contract] [nonce] [symbol] [ethereum-sender] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address] [flags]
ebcli query ethbridge prophecy 0x30753E4A8aad7F8597332E813735Def5dD395028 0 eth 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 --ethereum-chain-id=3 --token-contract-address=0x0000000000000000000000000000000000000000

# Confirm that the prophecy was successfully processed and that new token was minted to the testuser address
ebcli query account $(ebcli keys show testuser -a)

# Test out burning 1 of the eth for the return trip. We'll use "0x0000000000000000000000000000000000000000" for the token-contract-address, because we're dealing with the original EVM native token (eth).

# ebcli tx ethbridge burn [cosmos-sender-address] [ethereum-receiver-address] [amount] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address] [flags]
ebcli tx ethbridge burn $(ebcli keys show testuser -a) 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 1eth --ethereum-chain-id=3 --token-contract-address=0x0000000000000000000000000000000000000000 --from=testuser --yes

# Confirm that the token was successfully burned
ebcli query account $(ebcli keys show testuser -a)

# Test out locking up a cosmos stake coin for relaying over to the EVM chain. We'll use
# "0x345cA3e014Aaf5dcA488057592ee47305D9B3e10" for the token-contract-address because this is the addresses
# generated for the BridgeToken when following ./setup-eth-local.md.

# **_NOTE_** Make sure that you transferred some stake to the testuser from the validator account like described in one of the first instructions above, otherwise testuser will have insufficient funds to complete the transaction.

# ebcli tx ethbridge lock [cosmos-sender-address] [ethereum-receiver-address] [amount] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address] [flags]
ebcli tx ethbridge lock $(ebcli keys show testuser -a) 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 1stake  --ethereum-chain-id=3 --token-contract-address=0x345cA3e014Aaf5dcA488057592ee47305D9B3e10 --from=testuser --yes

# Confirm that the token was successfully locked
ebcli query account $(ebcli keys show testuser -a)

# Test out creating a bridge burn claim for the return trip back. This is similar to the create-claim we did earlier except for the asset being locked on the eth side, it was burned because the asset originated on the cosmos chain. Make sure you increment the nonce by one, since the first create-claim used nonce 0 this one should use nonce 1.

# ebcli tx ethbridge create-claim [bridge-registry-contract] [nonce] [symbol] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount] [claim-type] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address] [flags]
ebcli tx ethbridge create-claim 0x30753E4A8aad7F8597332E813735Def5dD395028 1 stake 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 1stake burn --ethereum-chain-id=3 --token-contract-address=0x345cA3e014Aaf5dcA488057592ee47305D9B3e10 --from=validator --yes

# Then read the prophecy to confirm it was created with the claim added
# ebcli query ethbridge prophecy [bridge-registry-contract] [nonce] [symbol] [ethereum-sender] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address] [flags]
ebcli query ethbridge prophecy 0x30753E4A8aad7F8597332E813735Def5dD395028 1 stake 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 --ethereum-chain-id=3 --token-contract-address=0x345cA3e014Aaf5dcA488057592ee47305D9B3e10

# Confirm that the prophecy was successfully processed and that stake coin was returned to the testuser address
ebcli query account $(ebcli keys show testuser -a)
```

To set up the EVM chain go to (the next step)[./setup-eth-local.md].
