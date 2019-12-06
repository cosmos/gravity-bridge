## Initialization

First, initialize a chain and create accounts:

```bash
# Initialize the genesis.json file that will help you to bootstrap the network
ebd init local --chain-id=peggy

# Create a key to hold your validator account and for another test account
ebcli keys add validator
# Enter password

ebcli keys add testuser
# Enter password

# Initialize the genesis account and transaction
ebd add-genesis-account $(ebcli keys show validator -a) 1000000000stake,1000000000atom

# Create genesis transaction
ebd gentx --name validator
# Enter password

# Collect genesis transaction
ebd collect-gentxs

# Now its safe to start `ebd`
ebd start
```

## Testing the application

Once you've initialized the application and started the Bridge blockchain with `ebd start`, you can test the available cli commands. They include sending tokens between accounts, querying accounts, claim creation, token burning, and token locking. Once the Relayer is running, you'll be able to submit new burning/locking txs to the chain using these commands.   

First, we'll test sending a random token in another terminal window.   

```bash
# Wait 10 seconds and in another terminal window, send 10 stake tokens from the validator to the testuser
ebcli tx send validator $(ebcli keys show testuser -a) 10stake --chain-id=peggy --yes

# Wait a few seconds for confirmation, then confirm token balances have changed appropriately
ebcli query account $(ebcli keys show validator -a) --trust-node
ebcli query account $(ebcli keys show testuser -a) --trust-node

# Then wait 10 seconds then confirm your validator was created correctly, and has become Bonded status
ebcli query staking validators --trust-node

# See the help for the ethbridge create claim function
ebcli tx ethbridge create-claim --help

# Now you can test out the ethbridge module by submitting a claim for an ethereum prophecy
# Create a bridge lock claim (Ethereum prophecies are stored on the blockchain with an identifier created by concatenating the nonce and sender address)
# ebcli tx ethbridge create-claim [ethereum-chain-id] [bridge-contract] [nonce] [symbol] [token-contract] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount]",
ebcli tx ethbridge create-claim 3 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB 0 eth 0x0000000000000000000000000000000000000000 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 3eth lock --from=validator --chain-id=peggy --yes

# Then read the prophecy to confirm it was created with the claim added
ebcli query ethbridge prophecy 3 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB 0 eth 0x0000000000000000000000000000000000000000 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 --trust-node

# Confirm that the prophecy was successfully processed and that new token was minted to the testuser address
ebcli query account $(ebcli keys show testuser -a) --trust-node

# Test out burning 1 of the eth for the return trip
ebcli tx ethbridge burn $(ebcli keys show testuser -a) 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 1eth --token-contract-address=0x682c2ae4053eac64cf1baaa04c739703dc043f0a --ethereum-chain-id=3 --from=testuser --chain-id=peggy --yes

# Confirm that the token was successfully burned
ebcli query account $(ebcli keys show testuser -a) --trust-node

# Test out locking up a cosmos stake coin for relaying over to ethereum
ebcli tx ethbridge lock $(ebcli keys show testuser -a) 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 1eth --token-contract-address=0x682c2ae4053eac64cf1baaa04c739703dc043f0a --ethereum-chain-id=3 --from=testuser --chain-id=peggy --yes

# Confirm that the token was successfully locked
ebcli query account $(ebcli keys show testuser -a) --trust-node

# Test out creating a bridge burn claim for the return trip back
ebcli tx ethbridge create-claim 1 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB 0 stake 0x3f5dab653144958ff6d309647baf1abde8da204d 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 1stake burn --from=validator --chain-id=peggy --yes

# Confirm that the prophecy was successfully processed and that stake coin was returned to the testuser address
ebcli query account $(ebcli keys show testuser -a) --trust-node
```