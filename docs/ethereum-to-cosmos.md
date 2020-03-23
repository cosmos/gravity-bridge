## Ethereum -> Cosmos asset transfers

### Sending Ethereum assets to Cosmos via Lock

Now we can lock our funds on the contracts by sending a lock transaction containing Eth/ERC20 assets. First, we'll use default parameters to lock Eth assets.  

Default parameter values for lock transactions:

- [COSMOS_RECIPIENT_ADDRESS] = `cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh`
- [TOKEN_CONTRACT_ADDRESS] = `eth` (Ethereum has no token contract and is denoted by 'eth')
- [WEI_AMOUNT] = `10`

```bash
# Open a new terminal window

# Send lock transaction with default parameters
yarn peggy:lock --default

# Send lock transaction with custom parameters
yarn peggy:lock [COSMOS_RECIPIENT_ADDRESS] [TOKEN_CONTRACT_ADDRESS] [WEI_AMOUNT]
```

`yarn peggy:lock --default` expected output in Relayer console (terminal 4):

```bash
2020/03/22 11:46:08 Witnessed tx 0xb799b5ed8df5f66c355b34fbcdbd132d0a0927c320c9b9c5ff7ea058ca55033c on block 16
2020/03/22 11:46:08 
Chain ID: 5777
Bridge contract address: 0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4
Token symbol: ETH
Token contract address: 0x0000000000000000000000000000000000000000
Sender: 0x627306090abaB3A6e1400e9345bC60c78a8BEf57
Recipient: cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh
Value: 10
Nonce: 3

{
  "height": "0",
  "txhash": "37C20C533345FC0A187A706E957BF509D56529F1B0FEA19AFC4278B23B2724A1",
  "raw_log": "[]"
}
```

### Testing ERC20 token support

The bridge supports the transfer of ERC20 token assets. First, we'll deploy a sample TEST token to the network and then use it to test the feature.

```bash
# Mint 1,000 TEST tokens to your account for local use
yarn token:mint

# Approve 100 TEST tokens to the Bridge contract
yarn token:approve --default

# You can also approve a custom amount of TEST tokens to the Bridge contract:
yarn token:approve 11

# Get deployed TEST token contract address
yarn token:address

# Lock TEST tokens on the Bridge contract
# Note: ERC20 token locking requires 3 custom params and does not support the --default flag
yarn peggy:lock [COSMOS_RECIPIENT_ADDRESS] [TEST_TOKEN_CONTRACT_ADDRESS] [TOKEN_AMOUNT]
```

`yarn peggy:lock` ERC20 expected output in ebrelayer console (with a `TOKEN_AMOUNT` of 11):

```bash
2020/03/22 11:48:09 Witnessed tx 0xab84de6d2f6bde3f2249cc1c31e23901432fa75b83a5b5b52c19e99479a797f1 on block 28
2020/03/22 11:48:09 
Chain ID: 5777
Bridge contract address: 0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4
Token symbol: TEST
Token contract address: 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB
Sender: 0x115F6e2004D7b4ccd6b9D5ab34e30909e0F612CD
Recipient: cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh
Value: 11
Nonce: 12

{
  "height": "0",
  "txhash": "013B79C59828872BA477FC8C2B98C155A0F8D520C42693363B7156F56B6C0A32",
  "raw_log": "[]"
}
```