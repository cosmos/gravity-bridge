package contract

// import (
//   "context"
//   "crypto/ecdsa"
//   "fmt"
//   "log"
//   "math/big"

//   "github.com/ethereum/go-ethereum/accounts/abi/bind"
//   "github.com/ethereum/go-ethereum/common"
//   "github.com/ethereum/go-ethereum/crypto"
//   "github.com/ethereum/go-ethereum/ethclient"

//   // Peggy "cmd/ebrelayer/contract/Peggy.sol"
// )

// func main() {
//   client, err := ethclient.Dial("https://ropsten.infura.io")
//   if err != nil {
//     log.Fatal(err)
//   }

//   privateKey, err := crypto.HexToECDSA("") //PK from config
//   if err != nil {
//     log.Fatal(err)
//   }

//   publicKey := privateKey.Public()
//   publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
//   if !ok {
//     log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
//   }

//   fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
//   nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
//   if err != nil {
//     log.Fatal(err)
//   }

//   gasPrice, err := client.SuggestGasPrice(context.Background())
//   if err != nil {
//     log.Fatal(err)
//   }

//   auth := bind.NewKeyedTransactor(privateKey)
//   auth.Nonce = big.NewInt(int64(nonce))
//   auth.Value = big.NewInt(20)     // in wei
//   auth.GasLimit = uint64(300000) // in units
//   auth.GasPrice = gasPrice

//   address := common.HexToAddress("0x3de4ef81Ba6243A60B0a32d3BCeD4173b6EA02bb")
//   instance, err := Peggy.lock("0x636f736d6f7331706a74677530766175326d35326e72796b64707a74727438383761796b756530687137646668", "0x0000000000000000000000000000000000000000", 20)
//   if err != nil {
//     log.Fatal(err)
//   }

//   key := [32]byte{}
//   value := [32]byte{}
//   copy(key[:], []byte("foo"))
//   copy(value[:], []byte("bar"))

//   tx, err := instance.SetItem(auth, key, value)
//   if err != nil {
//     log.Fatal(err)
//   }

//   fmt.Printf("tx sent: %s", tx.Hash().Hex())

//   result, err := instance.Items(nil, key)
//   if err != nil {
//     log.Fatal(err)
//   }

//   fmt.Println(string(result[:]))
// }