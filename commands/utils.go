package commands

import (
    "fmt"

    abci "github.com/tendermint/abci/types"
    "github.com/tendermint/basecoin/types"
    "github.com/tendermint/tendermint/rpc/client"
    "github.com/tendermint/go-wire"
    basecmd "github.com/tendermint/basecoin/cmd/basecoin/commands"
    tmtypes "github.com/tendermint/tendermint/types"
    
    "../plugins/etgate"
)


func GetAccWithClient(httpClient client.Client, address []byte) (*types.Account, error) {

	key := types.AccountKey(address)
	response, err := QueryWithClient(httpClient, key)
	if err != nil {
		return nil, err
	}

	accountBytes := response.Value

	if len(accountBytes) == 0 {
		return nil, fmt.Errorf("Account bytes are empty for address: %X ", address) //never stack trace
	}

	var acc *types.Account
	err = wire.ReadBinaryBytes(accountBytes, &acc)
	if err != nil {
		return nil, fmt.Errorf("Error reading account %X error: %v",
			accountBytes, err.Error())
	}

	return acc, nil
}

func QueryWithClient(httpClient client.Client, key []byte) (*abci.ResultQuery, error) {
	res, err := httpClient.ABCIQuery("/key", key, true)
	if err != nil {
		return nil, fmt.Errorf("Error calling /abci_query: %v", err)
	}
	if !res.Code.IsOK() {
		return nil, fmt.Errorf("Query got non-zero exit code: %v. %s", res.Code, res.Log)
	}
	return res.ResultQuery, nil
}

func BroadcastTxWithClient(httpClient client.Client, tx tmtypes.Tx) ([]byte, string, error) {
    res, err := httpClient.BroadcastTxCommit(tx)
    if err != nil {
        return nil, "", fmt.Errorf("Error on broadcast tx: %v", err)
    }

    if !res.CheckTx.Code.IsOK() {
        r := res.CheckTx
        return nil, "", fmt.Errorf("BroadcastTxCommit got non-zero exit code: %v, %X; %s", r.Code, r.Data, r.Log)
   }

    if !res.DeliverTx.Code.IsOK() {
        r := res.DeliverTx
        return nil, "", fmt.Errorf("BroadcastTxCommit got non-zero exit code: %v, %X; %s", r.Code, r.Data, r.Log)
    }

    return res.DeliverTx.Data, res.DeliverTx.Log, nil
}

func AppTx(httpClient client.Client, key *basecmd.Key, etgateTx etgate.ETGateTx, chainID string) error {
    acc, err := GetAccWithClient(httpClient, key.Address[:])
    if err != nil {
        return err
    }
    sequence := acc.Sequence + 1

    data := []byte(wire.BinaryBytes(struct {
        etgate.ETGateTx `json:"unwrap"`
    }{etgateTx}))
    
    smallCoins := types.Coin{Denom: "mycoin", Amount: 1}

    input := types.NewTxInput(key.PubKey, types.Coins{smallCoins}, sequence)
    tx := &types.AppTx {
        Gas: 0,
        Fee: smallCoins,
        Name: "ETGATE",
        Input: input,
        Data: data,
    }
    tx.Input.Signature = key.Sign(tx.SignBytes(chainID))
    txBytes := []byte(wire.BinaryBytes(struct {
        types.Tx `json:"unwrap"`
    }{tx}))

    data, log, err := BroadcastTxWithClient(httpClient, txBytes)
    if err != nil {
        return err
    }

    _, _ = data, log
    return nil
}
