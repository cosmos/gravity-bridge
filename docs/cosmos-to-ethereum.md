## Cosmos -> Ethereum asset transfers

### Sending Cosmos assets to Ethereum via Lock

At this point you should have a running Cosmos SDK application, a running EVM chain and a running relayer. If any of those are missing from your setup please go back to the [README](../README.md) and setup whatever is missing.

To send Cosmos assets to an EVM based chain, you'll use a transaction containing a lock message:

```bash
# In a new terminal window, send tokens to the testuser (10stake tokens)
ebcli tx send validator $(ebcli keys show testuser -a) 10stake --yes

# You can confirm they were received by querying the account
ebcli q account $(ebcli keys show testuser -a)
```

Now we can send the lock transaction for 1stake token to EVM chain address `0x5AEDA56215b167893e80B4fE645BA6d5Bab767DE`, which is the accounts[9] address of the truffle devlop local EVM chain.

```bash
# ebcli tx ethbridge lock [cosmos-sender-address] [ethereum-receiver-address] [amount] --ethereum-chain-id [ethereum-chain-id] [flags]
ebcli tx ethbridge lock $(ebcli keys show testuser -a) 0x5AEDA56215b167893e80B4fE645BA6d5Bab767DE 1 stake --ethereum-chain-id=3 --from=testuser --yes

```

Expected terminal output:

```bash
I[2020-04-23|23:33:13.092] New transaction witnessed                    
I[2020-04-23|23:33:13.092] 
Claim Type: lock
Cosmos Sender: cosmos1vnt63c0wtag5jnr6e9c7jz857amxrxcel0eucl
Ethereum Recipient: 0x5AEDA56215b167893e80B4fE645BA6d5Bab767DE
Symbol: STAKE
Amount: 1

Fetching CosmosBridge contract...
Sending new ProphecyClaim to CosmosBridge...
NewProphecyClaim tx hash: 0xe68c1fd54536e89bc28fd1803651e7184839ee3b8793a1cafe27f92212303e68
I[2020-04-23|23:33:13.198] Witnessed tx 0xe68c1fd54536e89bc28fd1803651e7184839ee3b8793a1cafe27f92212303e68 on block 21
 
I[2020-04-23|23:33:13.198] 
Prophecy ID: 4
Claim Type: 2
Sender: cosmos1vnt63c0wtag5jnr6e9c7jz857amxrxcel0eucl
Recipient: 0x5AEDA56215b167893e80B4fE645BA6d5Bab767DE
Symbol: PEGGYSTAKE
Token: 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
Amount: 1
Validator: 0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef
 
Generating unique message for ProphecyClaim 4
Signing message...
Signature generated: 0x646a7eb97d5cfc5171f5358c96c3b8e30e274f1ebbb604f50bc6b6d8f32bc60c6163e47555f7e48e4f3574b1358d81fe20fe7d2d6fbd47c635e5433b5ea2ed3b01
Tx Status: 1 - Successful

Fetching Oracle contract...
Sending new OracleClaim to Oracle...
NewOracleClaim tx hash: 0x5e53916ac1a8ce50564f97377c137417b49722183e4358ac6dee72b03c3af00d
Tx Status: 1 - Successful
```

To check the EVM chain balance of the account that just received the stake token you can use the following command. In our case we'd want to use the user address `0x5AEDA56215b167893e80B4fE645BA6d5Bab767DE` and the token address `0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA` which were just used in the last step. To check the EVM chain native asset balance just leave the token address blank.
```bash
yarn peggy:getTokenBalance [ACCOUNT_ADDRESS] [TOKEN_ADDRESS]
```

Congratulations, you've automatically relayed information from the lock transaction on the Cosmos SDK application to the contracts deployed on the Ethereum network as a new prophecy claim, witnessed the new prophecy claim, and signed its information to create an oracle claim. When a quorum of validators submit oracle claims for the prophecy, it will be processed. When a prophecy claim is successfully processed, the amount of tokens specified are minted by the contracts to the intended recipient on the Ethereum network.

### Returning Cosmos assets originally based on Ethereum via Burn

In the [Ethereum -> Cosmos asset transfers](./ethereum-to-cosmos.md) section, you sent assets to a Cosmos-SDK enabled chain. In order to return these assets to Ethereum and unlock the funds currently locked on the deployed contracts, you'll need to use a second type of transaction - `burn`. It's simple, just replace the ebcli `lock` command with `burn`:

To make sure you have EVM native eth on a cosmos account you can use first move some from the EVM chain with the following command:

```bash
yarn peggy:lock $(ebcli keys show testuser -a) eth 10
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

Now you can move the EVM native asset back to that vanity address (`0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9`) on the EVM chain.

```bash
# ebcli tx ethbridge burn [cosmos-sender-address] [ethereum-receiver-address] [amount] --ethereum-chain-id [ethereum-chain-id [flags]
ebcli tx ethbridge burn $(ebcli keys show testuser -a) 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 1eth  --ethereum-chain-id 3 --from testuser --yes
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
Reached threshold:       false
----------------------------------------
```   

Once the prophecy claim has reached the signed power threshold, it will be automatically processed and funds will be delivered to the intended Ethereum recipient.