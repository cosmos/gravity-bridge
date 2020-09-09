package types

type Valset struct {
	Nonce        int64
	Powers       []int64
	EthAddresses []string
}

func (v Valset) GetCheckpoint() []byte {
	// Getting the equivalent of solidity's abi.encodePacked (or abi.encode) does not seem to be straightforward
	// and I am skipping it for now to focus on the overall module structure
	// https://stackoverflow.com/questions/50772811/how-can-i-get-the-same-return-value-as-solidity-abi-encodepacked-in-golang
	return []byte("dothislater")
}
