## Cosmos -> Ethereum asset transfers

### Sending Cosmos assets to Ethereum via Lock

At this point you should have a running Cosmos SDK application, a running EVM chain and a running relayer. If any of those are missing from your setup please go back to the [README](../README.md) and setup whatever is missing.

To send Cosmos assets to an EVM based chain, you'll use a transaction containing a lock message:

```bash
# In a new terminal window, send tokens to the testuser (10stake tokens)
ebcli tx send validator $(ebcli keys show testuser -a) 10stake --yes

# You can confirm they were received by querying the account
ebcli q account $(ebcli keys show testuser -a)

# Before locking an asset on the Cosmos SDK application side, we need to also deploy a token on the EVM chain 
# that will represent the cosmos asset. To do this we use the following EVM command with the token name we'd like 
# to use (in this case "stake"):

yarn peggy:addBridgeToken stake

```
This should result in the following logs that contains the newly deployed `BridgeToken` contract address and adds it to the EVM token whitelist.

```bash
yarn peggy:addBridgeToken stake
yarn run v1.22.4
$ yarn workspace testnet-contracts peggy:addBridgeToken stake
$ truffle exec scripts/sendAddBridgeToken.js stake
Using network 'ganache'.

Fetching BridgeBank contract...
Attempting to send createNewBridgeToken() tx with symbol: 'stake'...
from 0x627306090abaB3A6e1400e9345bC60c78a8BEf57
Should deploy to 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
Bridge Token "stake" created at address: 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
Done in 30.87s.
```

Now we can send the lock transaction for 1stake token to EVM chain address `0x627306090abaB3A6e1400e9345bC60c78a8BEf57`, which is the accounts[0] address of the truffle devlop local EVM chain. We also use the newly created `BridgeToken` address from the previous step for `token-contract-address`.

```bash
# ebcli tx ethbridge lock [cosmos-sender-address] [ethereum-receiver-address] [amount] --ethereum-chain-id [ethereum-chain-id] [flags]
ebcli tx ethbridge lock $(ebcli keys show testuser -a) 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 1stake --ethereum-chain-id 3 --from testuser --yes
```

Expected terminal output:

```bash
I[2020-03-22|18:07:01.417] New transaction witnessed                    
I[2020-03-22|18:07:01.417] 
Claim Type: lock
Cosmos Sender: cosmos1vnt63c0wtag5jnr6e9c7jz857amxrxcel0eucl
Ethereum Recipient: 0x627306090abaB3A6e1400e9345bC60c78a8BEf57
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
Recipient: 0x627306090abaB3A6e1400e9345bC60c78a8BEf57
Symbol stake
Token 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
Amount: 1
Validator: 0xf17f52151EbEF6C7334FAD080c5704D77216b732

Generating unique message for ProphecyClaim 1
Signing message...
Signature generated: 0xae9b9ee377d85945d6516afc39c4d1f8efc1ead78ba03851d8e25cbf3227e3166a655c5cd280af9ff1a4f81ad501d754ae23b694ca834c45e0206d80504cd47b01
Tx Status: 1 - Successful

Fetching Oracle contract...
Sending new OracleClaim to Oracle...
NewOracleClaim tx hash: 0x89c1c905f65170e799fc17b16406aad61e07c857f3379190829f5fd5f9a157d9
Tx Status: 1 - Successful
```

To check the EVM chain balance of the account that just received the stake token you can use the following command. In our case we'd want to use the user address `0x627306090abaB3A6e1400e9345bC60c78a8BEf57` and the token address `0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA` which were just used in the last step. To check the EVM chain native asset balance just leave the token address blank.
```bash
yarn peggy:getTokenBalance [ACCOUNT_ADDRESS] [TOKEN_ADDRESS]
```

Congratulations, you've automatically relayed information from the lock transaction on the Cosmos SDK application to the contracts deployed on the Ethereum network as a new prophecy claim, witnessed the new prophecy claim, and signed its information to create an oracle claim. When a quorum of validators submit oracle claims for the prophecy, it will be processed. When a prophecy claim is successfully processed, the amount of tokens specified are minted by the contracts to the intended recipient on the Ethereum network.

### Returning Cosmos assets originally based on Ethereum via Burn

In the [Ethereum -> Cosmos asset transfers](./ethereum-to-cosmos.md) section, you sent assets to a Cosmos-SDK enabled chain. In order to return these assets to Ethereum and unlock the funds currently locked on the deployed contracts, you'll need to use a second type of transaction - `burn`. It's simple, just replace the ebcli `lock` command with `burn`:

To make sure you have EVM native eth on a cosmos account you can use first move some from the EVM chain with the following command:

```bash
yarn peggy:lock $(ebcli keys show testuser -a) eth 100000000000                                       
```
You can confirm this was successful using the following command:
```bash
ebcli q account $(ebcli keys show testuser -a)
```

Before moving the asset back to the EVM chain, check the balance of the destination address. It might be easier to see the token balance on an account that isn't being used to execute the transactions as well (which the `accounts[0]` address we previously used doing). For that you can try transferring it to the vanity address `0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9` with the following command:

```bash
yarn peggy:getTokenBalance 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 eth
```
You should see this is a balance of 0: `Eth balance for 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 is 0 Eth (0 Wei)`

Now you can move the EVM native asset back to that vanity address (`0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9`) on the EVM chain. Note we'll use token-contract-address `0x0000000000000000000000000000000000000000` since we're not actually transferring back a token but the EVM chain native asset (eth).

```bash
# ebcli tx ethbridge burn [cosmos-sender-address] [ethereum-receiver-address] [amount] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address] [flags]
ebcli tx ethbridge burn $(ebcli keys show testuser -a) 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 1eth  --ethereum-chain-id 3 --token-contract-address 0x0000000000000000000000000000000000000000 --from testuser --yes
```

You should now be able to see that address has received the ether:
```
> yarn peggy:getTokenBalance 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 eth

Eth balance for 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 is 0.000000000000000001 Eth (1 Wei)
```

<!-->
TODO: This isn't yet possible on the Ethereum contract side. Open issue https://github.com/cosmos/peggy/issues/123
### Returning EVM assets originally based on Cosmos via Burn

Similarly to the last step, it is possible to move Cosmos SDK application native assets which were previoulsy transferred to the eVM chain back to the Cosmos SDK application. Instead of locking the EVM based asset, as seen in the [Ethereum -> Cosmos](./ethereum-to-cosmos.md) document, we use the burn action instead. If you moved over some `stake` as described earlier in this document from the Cosmos SDK application to the EVM chain we can now move it back with a burn.

```
yarn peggy:burn [COSMOS_RECIPIENT_ADDRESS] [TOKEN_CONTRACT_ADDRESS] [WEI_AMOUNT]
```<!-->


### Prophecy claim processing

You are able to check the status of active prophecy claims. Prophecy claims can be processed once current signed power >= x% of total power, where x is the Oracle contract's consensus threshold parameter. This command is only for pending prophecy claims and will fail if the prophecy has already been confirmed.

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