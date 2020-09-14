### 1. Make tool

ở thư mục gốc của project chạy

```
make
```

### 2. Setup cosmos chain

```
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

để cho nó chạy mở 1 terminal khác

### 3. Deploy smart contract lên ropsten

Ở thư mục gốc của thư mục chạy

```
yarn peggy:all
```

Kiểm tra địa chỉ của các smart contract

```
cd testnet-contracts
truffle metwork
```

Kết quả trả về dạng

```
Network: ropsten (id: 3)
  BridgeBank: 0x32f82eeB46ed9c2AA0114e4B961cfcEBF18384Df
  BridgeRegistry: 0x14d268ed94340f757b253CdeBd3b3528B83aBdb1
  BridgeToken: 0x032A87fa8BA6031A4358213648B5eD5E72813A33
  CosmosBridge: 0x86C41cb7CbCC55919dFf52Ac9b8ac84D4ddBE2DA
  Migrations: 0x00c3b1ba15c5dD86Cf9253CA6b05e4617eeD3d3E
  Oracle: 0xad0023F9ebEF7741F399136F9Ca9276F3028b52B
  Valset: 0x1b8E9eBE7685D3d5a6f7Ce012f2e9738cD1E9Ef7
```

### 4. Chạy relayer

```
ebrelayer generate
ebrelayer init tcp://localhost:26657 wss://ropsten.infura.io/ws/v3/[Infura-Project-ID] [BridgeRegistry-ContractAddress] validator --chain-id=peggy
```

mở teminal khác lên

### 5. Chuyển tiền

#### 5.1 Chuyển ETH từ ethereum vào cosmos

Chạy lệnh sau để lấy address của cosmos receiver - tên là testuser

```
ebcli query account $(ebcli keys show testuser -a)
```

Kết quả trả về dạng

```
{
  "type": "cosmos-sdk/Account",
  "value": {
    "address": "cosmos1pgkwvwezfy3qkh99hjnf35ek3znzs79mwqf48y",
    "coins": [
      {
        "denom": "stake",
        "amount": "10"
      }
    ],
    "public_key": "cosmospub1addwnpepqwpznlktnvxvyxccslnp58janc6zk83huww6aynzq77ur2dvsfskct0atl9",
    "account_number": 3,
    "sequence": 7
  }
}
```

lấy giá trị của trường address thay vào đoạn

```
const DEFAULT_COSMOS_RECIPIENT = Web3.utils.utf8ToHex(
    'cosmos1pgkwvwezfy3qkh99hjnf35ek3znzs79mwqf48y'
  );
```

trong file testnet-contracts/scripts/sendLockTx.js

sau đó ở thư mục gốc chạy:

```
yarn token:lock --default
```

đợi 1 lúc query lại để thấy kết quả:

```
ebcli query account $(ebcli keys show testuser -a)
```

#### 5.2 Chuyển ETH từ cosmos về Ethereum

chạy lệnh sau

```
ebcli tx ethbridge burn $(ebcli keys show testuser -a) [địa chỉ ethereum nhận] 1000000000000000000 peggyeth --ethereum-chain-id=3 --from=testuser --yes
```
