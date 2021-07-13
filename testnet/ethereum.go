package main

type EthereumConfig struct {
	ChainID uint `json:"chainId"`
	HomesteadBlock uint `json:"homesteadBlock"`
	EIP150Block uint `json:"eip150Block"`
	EIP155Block uint `json:"eip155Block"`
	EIP158Block uint `json:"eip158Block"`
	ByzantiumBlock uint `json:"byzantiumBlock"`
	ConstantinopleBlock uint `json:"constantinopleBlock"`
	PetersburgBlock uint `json:"petersburgBlock"`
	IstanbulBlock uint `json:"istanbulBlock"`
	BerlinBlock uint `json:"berlinBlock"`
}

type Allocation struct {
	Balance string `json:"balance"`

}

type EthereumGenesis struct {
	Difficulty string                `json:"difficulty"`
	GasLimit   string                `json:"gasLimit"`
	Config     EthereumConfig        `json:"config"`
	Alloc      map[string]Allocation `json:"alloc"`
}
