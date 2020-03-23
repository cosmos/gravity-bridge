## Cosmos -> Ethereum asset transfers

### Sending Cosmos assets to Ethereum via Lock

To send Cosmos assets an EVM based chain, you'll use a transaction containing a lock message:

```bash
# In a new terminal window, send tokens to the testuser (10stake tokens)
ebcli tx send validator $(ebcli keys show testuser -a) 10stake --chain-id=peggy --yes

# Send lock transaction (1stake token)
ebcli tx ethbridge lock $(ebcli keys show testuser -a) [RECIPIENT_ETHEREUM_ADDRESS] 1stake --from testuser --chain-id peggy --ethereum-chain-id 3 --token-contract-address [TOKEN_CONTRACT_ADDRESS]
# Note: --token-contract-address will be '0x0000000000000000000000000000000000000000' for Ethereum assets
```

Expected terminal output:

```bash
I[2020-03-22|18:07:01.417] New transaction witnessed                    
I[2020-03-22|18:07:01.417] 
Claim Type: lock
Cosmos Sender: cosmos1vnt63c0wtag5jnr6e9c7jz857amxrxcel0eucl
Ethereum Recipient: 0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef
Token Address: 0x345cA3e014Aaf5dcA488057592ee47305D9B3e10
Symbol: stake
Amount: 1

Fetching CosmosBridge contract...
Sending new ProphecyClaim to CosmosBridge...
NewProphecyClaim tx hash: 0xd3ac2fb95e58e704c9c51bd171bb1b53623ae9505958105c86c09681bef46ec0
2020/03/22 18:07:01 Witnessed tx 0xd3ac2fb95e58e704c9c51bd171bb1b53623ae9505958105c86c09681bef46ec0 on block 17
2020/03/22 18:07:01 
Prophecy ID: 1
Claim Type: 2
Sender: cosmos1vnt63c0wtag5jnr6e9c7jz857amxrxcel0eucl
Recipient: 0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef
Symbol stake
Token 0x345cA3e014Aaf5dcA488057592ee47305D9B3e10
Amount: 1
Validator: 0xf17f52151EbEF6C7334FAD080c5704D77216b732

Generating unique message for ProphecyClaim 1
Signing message...
Signature generated: 0xc7ccfa125f92b5ec7780ce20948c4f5a174457bc3bfe025554507003fd42dcb67be0ea6e48c9a5493d2e63ea048f40ba81abd02945ac9ae8c69cc74409b2a14000
Tx Status: 1 - Successful

Fetching Oracle contract...
Sending new OracleClaim to Oracle...
NewOracleClaim tx hash: 0x89c1c905f65170e799fc17b16406aad61e07c857f3379190829f5fd5f9a157d9
Tx Status: 1 - Successful
```

Congratulations, you've automatically relayed information from the lock transaction on Tendermint to the contracts deployed on the Ethereum network as a new prophecy claim, witnessed the new prophecy claim, and signed its information to create an oracle claim. When enough validators submit oracle claims for the prophecy claim, it will be processed. When a prophecy claim is successfully processed, the amount of tokens specified will be minted by the contracts to the intended recipient on the Ethereum network.

### Returning Cosmos assets originally based on Ethereum via Burn

In the `Ethereum -> Cosmos asset transfers` section, you sent assets to a Cosmos-SDK enabled chain. In order to return these assets to Ethereum and unlock the funds currently locked on the deployed contracts, you'll need to use a second type of transaction - burn. It's simple, just replace the ebcli `lock` command with `burn`:

```bash
# Send burn transaction (1stake token)
ebcli tx ethbridge burn $(ebcli keys show testuser -a) [RECIPIENT_ETHEREUM_ADDRESS] 1stake --from testuser --chain-id peggy --ethereum-chain-id 3 --token-contract-address [TOKEN_CONTRACT_ADDRESS]
```

### Prophecy claim processing

You are able to check the status of active prophecy claims. Prophecy claims reach the signed power threshold when the signed power surpasses the total power threshold.

```bash
# Check prophecy claim status
yarn peggy:check [PROPHECY_CLAIM_ID]
```

Expected output (for a prophecy claim with an ID of 2)

```bash
Fetching Oracle contract...
Attempting to send checkBridgeProphecy() tx...

        Prophecy 2 status:
----------------------------------------
Weighted total power:    104
Weighted signed power:   150
Reached threshold:       true
----------------------------------------
```   

Once the prophecy claim has reached the signed power threshold, anyone may initiate its processing. Any attempts to process prophecy claims under the signed power threshold will be rejected by the contracts.   

```bash
# Process the prophecy claim
yarn peggy:process [PROPHECY_CLAIM_ID]
```