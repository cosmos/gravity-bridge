// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package CosmosBridge

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// CosmosBridgeABI is the input ABI used to generate the binding from.
const CosmosBridgeABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"bridgeBank\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_prophecyID\",\"type\":\"uint256\"}],\"name\":\"isProphecyClaimValidatorActive\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_nonce\",\"type\":\"uint256\"},{\"name\":\"_cosmosSender\",\"type\":\"bytes\"},{\"name\":\"_ethereumReceiver\",\"type\":\"address\"},{\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"newProphecyClaim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"hasBridgeBank\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_prophecyID\",\"type\":\"uint256\"}],\"name\":\"completeProphecyClaim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"valset\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_bridgeBank\",\"type\":\"address\"}],\"name\":\"setBridgeBank\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"prophecyClaimCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_prophecyID\",\"type\":\"uint256\"}],\"name\":\"isProphecyClaimActive\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"prophecyClaims\",\"outputs\":[{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"cosmosSender\",\"type\":\"bytes\"},{\"name\":\"ethereumReceiver\",\"type\":\"address\"},{\"name\":\"originalValidator\",\"type\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"symbol\",\"type\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"hasOracle\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"},{\"name\":\"_valset\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"LogOracleSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_bridgeBank\",\"type\":\"address\"}],\"name\":\"LogBridgeBankSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_cosmosSender\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"_ethereumReceiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_symbol\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"LogNewProphecyClaim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prophecyID\",\"type\":\"uint256\"}],\"name\":\"LogProphecyCompleted\",\"type\":\"event\"}]"

// CosmosBridgeBin is the compiled bytecode used for deploying new contracts.
const CosmosBridgeBin = `608060405234801561001057600080fd5b50604051604080611f968339810180604052604081101561003057600080fd5b8101908080519060200190929190805190602001909291905050506000600481905550816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600260146101000a81548160ff0219169083151502179055506000600360146101000a81548160ff0219169083151502179055505050611e7b8061011b6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80637dc0d1d01161008c5780638ea5352d116100665780638ea5352d146104db578063d8da69ea146104f9578063db4237af1461053f578063fb7831f21461070e576100ea565b80637dc0d1d0146104035780637f54af0c1461044d578063814c92c314610497576100ea565b80635a82e399116100c85780635a82e399146101c957806369294a4e1461036f5780636b3ce98c146103915780637adbf973146103bf576100ea565b80630e41f373146100ef578063529f3dd214610139578063570ca7351461017f575b600080fd5b6100f7610730565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101656004803603602081101561014f57600080fd5b8101908080359060200190929190505050610756565b604051808215151515815260200191505060405180910390f35b61018761086f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61036d600480360360c08110156101df57600080fd5b81019080803590602001909291908035906020019064010000000081111561020657600080fd5b82018360208201111561021857600080fd5b8035906020019184600183028401116401000000008311171561023a57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001906401000000008111156102dd57600080fd5b8201836020820111156102ef57600080fd5b8035906020019184600183028401116401000000008311171561031157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929080359060200190929190505050610894565b005b610377610e4c565b604051808215151515815260200191505060405180910390f35b6103bd600480360360208110156103a757600080fd5b8101908080359060200190929190505050610e5f565b005b610401600480360360208110156103d557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ffd565b005b61040b611209565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61045561122f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6104d9600480360360208110156104ad57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611255565b005b6104e3611461565b6040518082815260200191505060405180910390f35b6105256004803603602081101561050f57600080fd5b8101908080359060200190929190505050611467565b604051808215151515815260200191505060405180910390f35b61056b6004803603602081101561055557600080fd5b81019080803590602001909291905050506114ad565b60405180898152602001806020018873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018060200185815260200184600381111561062557fe5b60ff16815260200183810383528a818151815260200191508051906020019080838360005b8381101561066557808201518184015260208101905061064a565b50505050905090810190601f1680156106925780820380516001836020036101000a031916815260200191505b50838103825286818151815260200191508051906020019080838360005b838110156106cb5780820151818401526020810190506106b0565b50505050905090810190601f1680156106f85780820380516001836020036101000a031916815260200191505b509a505050505050505050505060405180910390f35b610716611692565b604051808215151515815260200191505060405180910390f35b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340550a1c6005600085815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff166040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561082d57600080fd5b505afa158015610841573d6000803e3d6000fd5b505050506040513d602081101561085757600080fd5b81019080805190602001909291905050509050919050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60011515600260149054906101000a900460ff1615151480156108ca575060011515600360149054906101000a900460ff161515145b61091f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526046815260200180611d7c6046913960600191505060405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340550a1c336040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156109be57600080fd5b505afa1580156109d2573d6000803e3d6000fd5b505050506040513d60208110156109e857600080fd5b8101908080519060200190929190505050610a6b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f4d75737420626520616e206163746976652076616c696461746f72000000000081525060200191505060405180910390fd5b610a8160016004546116a590919063ffffffff16565b6004819055506000339050610a94611bc4565b6040518061010001604052808981526020018881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff16815260200185815260200184815260200160016003811115610b1957fe5b815250905080600560006004548152602001908152602001600020600082015181600001556020820151816001019080519060200190610b5a929190611c56565b5060408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060608201518160030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160040160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a0820151816005019080519060200190610c4c929190611cd6565b5060c0820151816006015560e08201518160070160006101000a81548160ff02191690836003811115610c7b57fe5b02179055509050507f48597fd2bac5835331b7653bb9b6284df26bb7fbd712b0059a42baf36e1b15e9600454898989868a8a8a60405180898152602001888152602001806020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001848152602001838103835289818151815260200191508051906020019080838360005b83811015610d9a578082015181840152602081019050610d7f565b50505050905090810190601f168015610dc75780820380516001836020036101000a031916815260200191505b50838103825285818151815260200191508051906020019080838360005b83811015610e00578082015181840152602081019050610de5565b50505050905090810190601f168015610e2d5780820380516001836020036101000a031916815260200191505b509a505050505050505050505060405180910390a15050505050505050565b600360149054906101000a900460ff1681565b80610e6981611467565b610edb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f50726f706865637920636c61696d206973206e6f74206163746976650000000081525060200191505060405180910390fd5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610f81576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526027815260200180611e296027913960400191505060405180910390fd5b60026005600084815260200190815260200160002060070160006101000a81548160ff02191690836003811115610fb457fe5b0217905550610fc28261172d565b7f18083372d5d670784eb4cddbb916a6b39bae1015476bdfe19a0fb1482c45313d826040518082815260200191505060405180910390a15050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146110bf576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4d75737420626520746865206f70657261746f722e000000000000000000000081525060200191505060405180910390fd5b600260149054906101000a900460ff1615611125576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526031815260200180611df86031913960400191505060405180910390fd5b6001600260146101000a81548160ff02191690831515021790555080600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f6efb0434342713e2e9b1501dbebf76b4ed18406ea77ab5d56535cc26dec3adc0600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a150565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611317576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4d75737420626520746865206f70657261746f722e000000000000000000000081525060200191505060405180910390fd5b600360149054906101000a900460ff161561137d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526036815260200180611dc26036913960400191505060405180910390fd5b6001600360146101000a81548160ff02191690831515021790555080600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507fc8b65043fb196ac032b79a435397d1d14a96b4e9d12e366c3b1f550cb01d2dfa600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a150565b60045481565b60006001600381111561147657fe5b6005600084815260200190815260200160002060070160009054906101000a900460ff1660038111156114a557fe5b149050919050565b6005602052806000526040600020600091509050806000015490806001018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561155f5780601f106115345761010080835404028352916020019161155f565b820191906000526020600020905b81548152906001019060200180831161154257829003601f168201915b5050505050908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060040160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806005018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561166f5780601f106116445761010080835404028352916020019161166f565b820191906000526020600020905b81548152906001019060200180831161165257829003601f168201915b5050505050908060060154908060070160009054906101000a900460ff16905088565b600260149054906101000a900460ff1681565b600080828401905083811015611723576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b611735611bc4565b600560008381526020019081526020016000206040518061010001604052908160008201548152602001600182018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156117f75780601f106117cc576101008083540402835291602001916117f7565b820191906000526020600020905b8154815290600101906020018083116117da57829003601f168201915b505050505081526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016003820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016004820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600582018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561199b5780601f106119705761010080835404028352916020019161199b565b820191906000526020600020905b81548152906001019060200180831161197e57829003601f168201915b50505050508152602001600682015481526020016007820160009054906101000a900460ff1660038111156119cc57fe5b60038111156119d757fe5b815250509050600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663cdf68c418260200151836040015184608001518560a001518660c001516040518663ffffffff1660e01b815260040180806020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001848152602001838103835288818151815260200191508051906020019080838360005b83811015611af1578082015181840152602081019050611ad6565b50505050905090810190601f168015611b1e5780820380516001836020036101000a031916815260200191505b50838103825285818151815260200191508051906020019080838360005b83811015611b57578082015181840152602081019050611b3c565b50505050905090810190601f168015611b845780820380516001836020036101000a031916815260200191505b50975050505050505050600060405180830381600087803b158015611ba857600080fd5b505af1158015611bbc573d6000803e3d6000fd5b505050505050565b6040518061010001604052806000815260200160608152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001606081526020016000815260200160006003811115611c5057fe5b81525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611c9757805160ff1916838001178555611cc5565b82800160010185558215611cc5579182015b82811115611cc4578251825591602001919060010190611ca9565b5b509050611cd29190611d56565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611d1757805160ff1916838001178555611d45565b82800160010185558215611d45579182015b82811115611d44578251825591602001919060010190611d29565b5b509050611d529190611d56565b5090565b611d7891905b80821115611d74576000816000905550600101611d5c565b5090565b9056fe546865204f70657261746f72206d7573742073657420746865206f7261636c6520616e64206272696467652062616e6b20666f72206272696467652061637469766174696f6e546865204272696467652042616e6b2063616e6e6f742062652075706461746564206f6e636520697420686173206265656e20736574546865204f7261636c652063616e6e6f742062652075706461746564206f6e636520697420686173206265656e207365744f6e6c7920746865204f7261636c65206d617920636f6d706c6574652070726f70686563696573a165627a7a72305820ade5acbca5c22f0c94f1d4848d0a98e279d9bfabd99ae5c63e81be2328e358b70029`

// DeployCosmosBridge deploys a new Ethereum contract, binding an instance of CosmosBridge to it.
func DeployCosmosBridge(auth *bind.TransactOpts, backend bind.ContractBackend, _operator common.Address, _valset common.Address) (common.Address, *types.Transaction, *CosmosBridge, error) {
	parsed, err := abi.JSON(strings.NewReader(CosmosBridgeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CosmosBridgeBin), backend, _operator, _valset)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CosmosBridge{CosmosBridgeCaller: CosmosBridgeCaller{contract: contract}, CosmosBridgeTransactor: CosmosBridgeTransactor{contract: contract}, CosmosBridgeFilterer: CosmosBridgeFilterer{contract: contract}}, nil
}

// CosmosBridge is an auto generated Go binding around an Ethereum contract.
type CosmosBridge struct {
	CosmosBridgeCaller     // Read-only binding to the contract
	CosmosBridgeTransactor // Write-only binding to the contract
	CosmosBridgeFilterer   // Log filterer for contract events
}

// CosmosBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type CosmosBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CosmosBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CosmosBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CosmosBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CosmosBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CosmosBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CosmosBridgeSession struct {
	Contract     *CosmosBridge     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CosmosBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CosmosBridgeCallerSession struct {
	Contract *CosmosBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// CosmosBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CosmosBridgeTransactorSession struct {
	Contract     *CosmosBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// CosmosBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type CosmosBridgeRaw struct {
	Contract *CosmosBridge // Generic contract binding to access the raw methods on
}

// CosmosBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CosmosBridgeCallerRaw struct {
	Contract *CosmosBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// CosmosBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CosmosBridgeTransactorRaw struct {
	Contract *CosmosBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCosmosBridge creates a new instance of CosmosBridge, bound to a specific deployed contract.
func NewCosmosBridge(address common.Address, backend bind.ContractBackend) (*CosmosBridge, error) {
	contract, err := bindCosmosBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CosmosBridge{CosmosBridgeCaller: CosmosBridgeCaller{contract: contract}, CosmosBridgeTransactor: CosmosBridgeTransactor{contract: contract}, CosmosBridgeFilterer: CosmosBridgeFilterer{contract: contract}}, nil
}

// NewCosmosBridgeCaller creates a new read-only instance of CosmosBridge, bound to a specific deployed contract.
func NewCosmosBridgeCaller(address common.Address, caller bind.ContractCaller) (*CosmosBridgeCaller, error) {
	contract, err := bindCosmosBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeCaller{contract: contract}, nil
}

// NewCosmosBridgeTransactor creates a new write-only instance of CosmosBridge, bound to a specific deployed contract.
func NewCosmosBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*CosmosBridgeTransactor, error) {
	contract, err := bindCosmosBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeTransactor{contract: contract}, nil
}

// NewCosmosBridgeFilterer creates a new log filterer instance of CosmosBridge, bound to a specific deployed contract.
func NewCosmosBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*CosmosBridgeFilterer, error) {
	contract, err := bindCosmosBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeFilterer{contract: contract}, nil
}

// bindCosmosBridge binds a generic wrapper to an already deployed contract.
func bindCosmosBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CosmosBridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CosmosBridge *CosmosBridgeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CosmosBridge.Contract.CosmosBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CosmosBridge *CosmosBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CosmosBridge.Contract.CosmosBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CosmosBridge *CosmosBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CosmosBridge.Contract.CosmosBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CosmosBridge *CosmosBridgeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CosmosBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CosmosBridge *CosmosBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CosmosBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CosmosBridge *CosmosBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CosmosBridge.Contract.contract.Transact(opts, method, params...)
}

// BridgeBank is a free data retrieval call binding the contract method 0x0e41f373.
//
// Solidity: function bridgeBank() constant returns(address)
func (_CosmosBridge *CosmosBridgeCaller) BridgeBank(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "bridgeBank")
	return *ret0, err
}

// BridgeBank is a free data retrieval call binding the contract method 0x0e41f373.
//
// Solidity: function bridgeBank() constant returns(address)
func (_CosmosBridge *CosmosBridgeSession) BridgeBank() (common.Address, error) {
	return _CosmosBridge.Contract.BridgeBank(&_CosmosBridge.CallOpts)
}

// BridgeBank is a free data retrieval call binding the contract method 0x0e41f373.
//
// Solidity: function bridgeBank() constant returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) BridgeBank() (common.Address, error) {
	return _CosmosBridge.Contract.BridgeBank(&_CosmosBridge.CallOpts)
}

// HasBridgeBank is a free data retrieval call binding the contract method 0x69294a4e.
//
// Solidity: function hasBridgeBank() constant returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) HasBridgeBank(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "hasBridgeBank")
	return *ret0, err
}

// HasBridgeBank is a free data retrieval call binding the contract method 0x69294a4e.
//
// Solidity: function hasBridgeBank() constant returns(bool)
func (_CosmosBridge *CosmosBridgeSession) HasBridgeBank() (bool, error) {
	return _CosmosBridge.Contract.HasBridgeBank(&_CosmosBridge.CallOpts)
}

// HasBridgeBank is a free data retrieval call binding the contract method 0x69294a4e.
//
// Solidity: function hasBridgeBank() constant returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) HasBridgeBank() (bool, error) {
	return _CosmosBridge.Contract.HasBridgeBank(&_CosmosBridge.CallOpts)
}

// HasOracle is a free data retrieval call binding the contract method 0xfb7831f2.
//
// Solidity: function hasOracle() constant returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) HasOracle(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "hasOracle")
	return *ret0, err
}

// HasOracle is a free data retrieval call binding the contract method 0xfb7831f2.
//
// Solidity: function hasOracle() constant returns(bool)
func (_CosmosBridge *CosmosBridgeSession) HasOracle() (bool, error) {
	return _CosmosBridge.Contract.HasOracle(&_CosmosBridge.CallOpts)
}

// HasOracle is a free data retrieval call binding the contract method 0xfb7831f2.
//
// Solidity: function hasOracle() constant returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) HasOracle() (bool, error) {
	return _CosmosBridge.Contract.HasOracle(&_CosmosBridge.CallOpts)
}

// IsProphecyClaimActive is a free data retrieval call binding the contract method 0xd8da69ea.
//
// Solidity: function isProphecyClaimActive(uint256 _prophecyID) constant returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) IsProphecyClaimActive(opts *bind.CallOpts, _prophecyID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "isProphecyClaimActive", _prophecyID)
	return *ret0, err
}

// IsProphecyClaimActive is a free data retrieval call binding the contract method 0xd8da69ea.
//
// Solidity: function isProphecyClaimActive(uint256 _prophecyID) constant returns(bool)
func (_CosmosBridge *CosmosBridgeSession) IsProphecyClaimActive(_prophecyID *big.Int) (bool, error) {
	return _CosmosBridge.Contract.IsProphecyClaimActive(&_CosmosBridge.CallOpts, _prophecyID)
}

// IsProphecyClaimActive is a free data retrieval call binding the contract method 0xd8da69ea.
//
// Solidity: function isProphecyClaimActive(uint256 _prophecyID) constant returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) IsProphecyClaimActive(_prophecyID *big.Int) (bool, error) {
	return _CosmosBridge.Contract.IsProphecyClaimActive(&_CosmosBridge.CallOpts, _prophecyID)
}

// IsProphecyClaimValidatorActive is a free data retrieval call binding the contract method 0x529f3dd2.
//
// Solidity: function isProphecyClaimValidatorActive(uint256 _prophecyID) constant returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) IsProphecyClaimValidatorActive(opts *bind.CallOpts, _prophecyID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "isProphecyClaimValidatorActive", _prophecyID)
	return *ret0, err
}

// IsProphecyClaimValidatorActive is a free data retrieval call binding the contract method 0x529f3dd2.
//
// Solidity: function isProphecyClaimValidatorActive(uint256 _prophecyID) constant returns(bool)
func (_CosmosBridge *CosmosBridgeSession) IsProphecyClaimValidatorActive(_prophecyID *big.Int) (bool, error) {
	return _CosmosBridge.Contract.IsProphecyClaimValidatorActive(&_CosmosBridge.CallOpts, _prophecyID)
}

// IsProphecyClaimValidatorActive is a free data retrieval call binding the contract method 0x529f3dd2.
//
// Solidity: function isProphecyClaimValidatorActive(uint256 _prophecyID) constant returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) IsProphecyClaimValidatorActive(_prophecyID *big.Int) (bool, error) {
	return _CosmosBridge.Contract.IsProphecyClaimValidatorActive(&_CosmosBridge.CallOpts, _prophecyID)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_CosmosBridge *CosmosBridgeCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "operator")
	return *ret0, err
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_CosmosBridge *CosmosBridgeSession) Operator() (common.Address, error) {
	return _CosmosBridge.Contract.Operator(&_CosmosBridge.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) Operator() (common.Address, error) {
	return _CosmosBridge.Contract.Operator(&_CosmosBridge.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() constant returns(address)
func (_CosmosBridge *CosmosBridgeCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "oracle")
	return *ret0, err
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() constant returns(address)
func (_CosmosBridge *CosmosBridgeSession) Oracle() (common.Address, error) {
	return _CosmosBridge.Contract.Oracle(&_CosmosBridge.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() constant returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) Oracle() (common.Address, error) {
	return _CosmosBridge.Contract.Oracle(&_CosmosBridge.CallOpts)
}

// ProphecyClaimCount is a free data retrieval call binding the contract method 0x8ea5352d.
//
// Solidity: function prophecyClaimCount() constant returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) ProphecyClaimCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "prophecyClaimCount")
	return *ret0, err
}

// ProphecyClaimCount is a free data retrieval call binding the contract method 0x8ea5352d.
//
// Solidity: function prophecyClaimCount() constant returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) ProphecyClaimCount() (*big.Int, error) {
	return _CosmosBridge.Contract.ProphecyClaimCount(&_CosmosBridge.CallOpts)
}

// ProphecyClaimCount is a free data retrieval call binding the contract method 0x8ea5352d.
//
// Solidity: function prophecyClaimCount() constant returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) ProphecyClaimCount() (*big.Int, error) {
	return _CosmosBridge.Contract.ProphecyClaimCount(&_CosmosBridge.CallOpts)
}

// ProphecyClaims is a free data retrieval call binding the contract method 0xdb4237af.
//
// Solidity: function prophecyClaims(uint256 ) constant returns(uint256 nonce, bytes cosmosSender, address ethereumReceiver, address originalValidator, address tokenAddress, string symbol, uint256 amount, uint8 status)
func (_CosmosBridge *CosmosBridgeCaller) ProphecyClaims(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Nonce             *big.Int
	CosmosSender      []byte
	EthereumReceiver  common.Address
	OriginalValidator common.Address
	TokenAddress      common.Address
	Symbol            string
	Amount            *big.Int
	Status            uint8
}, error) {
	ret := new(struct {
		Nonce             *big.Int
		CosmosSender      []byte
		EthereumReceiver  common.Address
		OriginalValidator common.Address
		TokenAddress      common.Address
		Symbol            string
		Amount            *big.Int
		Status            uint8
	})
	out := ret
	err := _CosmosBridge.contract.Call(opts, out, "prophecyClaims", arg0)
	return *ret, err
}

// ProphecyClaims is a free data retrieval call binding the contract method 0xdb4237af.
//
// Solidity: function prophecyClaims(uint256 ) constant returns(uint256 nonce, bytes cosmosSender, address ethereumReceiver, address originalValidator, address tokenAddress, string symbol, uint256 amount, uint8 status)
func (_CosmosBridge *CosmosBridgeSession) ProphecyClaims(arg0 *big.Int) (struct {
	Nonce             *big.Int
	CosmosSender      []byte
	EthereumReceiver  common.Address
	OriginalValidator common.Address
	TokenAddress      common.Address
	Symbol            string
	Amount            *big.Int
	Status            uint8
}, error) {
	return _CosmosBridge.Contract.ProphecyClaims(&_CosmosBridge.CallOpts, arg0)
}

// ProphecyClaims is a free data retrieval call binding the contract method 0xdb4237af.
//
// Solidity: function prophecyClaims(uint256 ) constant returns(uint256 nonce, bytes cosmosSender, address ethereumReceiver, address originalValidator, address tokenAddress, string symbol, uint256 amount, uint8 status)
func (_CosmosBridge *CosmosBridgeCallerSession) ProphecyClaims(arg0 *big.Int) (struct {
	Nonce             *big.Int
	CosmosSender      []byte
	EthereumReceiver  common.Address
	OriginalValidator common.Address
	TokenAddress      common.Address
	Symbol            string
	Amount            *big.Int
	Status            uint8
}, error) {
	return _CosmosBridge.Contract.ProphecyClaims(&_CosmosBridge.CallOpts, arg0)
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() constant returns(address)
func (_CosmosBridge *CosmosBridgeCaller) Valset(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "valset")
	return *ret0, err
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() constant returns(address)
func (_CosmosBridge *CosmosBridgeSession) Valset() (common.Address, error) {
	return _CosmosBridge.Contract.Valset(&_CosmosBridge.CallOpts)
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() constant returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) Valset() (common.Address, error) {
	return _CosmosBridge.Contract.Valset(&_CosmosBridge.CallOpts)
}

// CompleteProphecyClaim is a paid mutator transaction binding the contract method 0x6b3ce98c.
//
// Solidity: function completeProphecyClaim(uint256 _prophecyID) returns()
func (_CosmosBridge *CosmosBridgeTransactor) CompleteProphecyClaim(opts *bind.TransactOpts, _prophecyID *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "completeProphecyClaim", _prophecyID)
}

// CompleteProphecyClaim is a paid mutator transaction binding the contract method 0x6b3ce98c.
//
// Solidity: function completeProphecyClaim(uint256 _prophecyID) returns()
func (_CosmosBridge *CosmosBridgeSession) CompleteProphecyClaim(_prophecyID *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.CompleteProphecyClaim(&_CosmosBridge.TransactOpts, _prophecyID)
}

// CompleteProphecyClaim is a paid mutator transaction binding the contract method 0x6b3ce98c.
//
// Solidity: function completeProphecyClaim(uint256 _prophecyID) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) CompleteProphecyClaim(_prophecyID *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.CompleteProphecyClaim(&_CosmosBridge.TransactOpts, _prophecyID)
}

// NewProphecyClaim is a paid mutator transaction binding the contract method 0x5a82e399.
//
// Solidity: function newProphecyClaim(uint256 _nonce, bytes _cosmosSender, address _ethereumReceiver, address _tokenAddress, string _symbol, uint256 _amount) returns()
func (_CosmosBridge *CosmosBridgeTransactor) NewProphecyClaim(opts *bind.TransactOpts, _nonce *big.Int, _cosmosSender []byte, _ethereumReceiver common.Address, _tokenAddress common.Address, _symbol string, _amount *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "newProphecyClaim", _nonce, _cosmosSender, _ethereumReceiver, _tokenAddress, _symbol, _amount)
}

// NewProphecyClaim is a paid mutator transaction binding the contract method 0x5a82e399.
//
// Solidity: function newProphecyClaim(uint256 _nonce, bytes _cosmosSender, address _ethereumReceiver, address _tokenAddress, string _symbol, uint256 _amount) returns()
func (_CosmosBridge *CosmosBridgeSession) NewProphecyClaim(_nonce *big.Int, _cosmosSender []byte, _ethereumReceiver common.Address, _tokenAddress common.Address, _symbol string, _amount *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.NewProphecyClaim(&_CosmosBridge.TransactOpts, _nonce, _cosmosSender, _ethereumReceiver, _tokenAddress, _symbol, _amount)
}

// NewProphecyClaim is a paid mutator transaction binding the contract method 0x5a82e399.
//
// Solidity: function newProphecyClaim(uint256 _nonce, bytes _cosmosSender, address _ethereumReceiver, address _tokenAddress, string _symbol, uint256 _amount) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) NewProphecyClaim(_nonce *big.Int, _cosmosSender []byte, _ethereumReceiver common.Address, _tokenAddress common.Address, _symbol string, _amount *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.NewProphecyClaim(&_CosmosBridge.TransactOpts, _nonce, _cosmosSender, _ethereumReceiver, _tokenAddress, _symbol, _amount)
}

// SetBridgeBank is a paid mutator transaction binding the contract method 0x814c92c3.
//
// Solidity: function setBridgeBank(address _bridgeBank) returns()
func (_CosmosBridge *CosmosBridgeTransactor) SetBridgeBank(opts *bind.TransactOpts, _bridgeBank common.Address) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "setBridgeBank", _bridgeBank)
}

// SetBridgeBank is a paid mutator transaction binding the contract method 0x814c92c3.
//
// Solidity: function setBridgeBank(address _bridgeBank) returns()
func (_CosmosBridge *CosmosBridgeSession) SetBridgeBank(_bridgeBank common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.SetBridgeBank(&_CosmosBridge.TransactOpts, _bridgeBank)
}

// SetBridgeBank is a paid mutator transaction binding the contract method 0x814c92c3.
//
// Solidity: function setBridgeBank(address _bridgeBank) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) SetBridgeBank(_bridgeBank common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.SetBridgeBank(&_CosmosBridge.TransactOpts, _bridgeBank)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_CosmosBridge *CosmosBridgeTransactor) SetOracle(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "setOracle", _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_CosmosBridge *CosmosBridgeSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.SetOracle(&_CosmosBridge.TransactOpts, _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.SetOracle(&_CosmosBridge.TransactOpts, _oracle)
}

// CosmosBridgeLogBridgeBankSetIterator is returned from FilterLogBridgeBankSet and is used to iterate over the raw logs and unpacked data for LogBridgeBankSet events raised by the CosmosBridge contract.
type CosmosBridgeLogBridgeBankSetIterator struct {
	Event *CosmosBridgeLogBridgeBankSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CosmosBridgeLogBridgeBankSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogBridgeBankSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CosmosBridgeLogBridgeBankSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CosmosBridgeLogBridgeBankSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogBridgeBankSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogBridgeBankSet represents a LogBridgeBankSet event raised by the CosmosBridge contract.
type CosmosBridgeLogBridgeBankSet struct {
	BridgeBank common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogBridgeBankSet is a free log retrieval operation binding the contract event 0xc8b65043fb196ac032b79a435397d1d14a96b4e9d12e366c3b1f550cb01d2dfa.
//
// Solidity: event LogBridgeBankSet(address _bridgeBank)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogBridgeBankSet(opts *bind.FilterOpts) (*CosmosBridgeLogBridgeBankSetIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogBridgeBankSet")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogBridgeBankSetIterator{contract: _CosmosBridge.contract, event: "LogBridgeBankSet", logs: logs, sub: sub}, nil
}

// WatchLogBridgeBankSet is a free log subscription operation binding the contract event 0xc8b65043fb196ac032b79a435397d1d14a96b4e9d12e366c3b1f550cb01d2dfa.
//
// Solidity: event LogBridgeBankSet(address _bridgeBank)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogBridgeBankSet(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogBridgeBankSet) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogBridgeBankSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogBridgeBankSet)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogBridgeBankSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// CosmosBridgeLogNewProphecyClaimIterator is returned from FilterLogNewProphecyClaim and is used to iterate over the raw logs and unpacked data for LogNewProphecyClaim events raised by the CosmosBridge contract.
type CosmosBridgeLogNewProphecyClaimIterator struct {
	Event *CosmosBridgeLogNewProphecyClaim // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CosmosBridgeLogNewProphecyClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogNewProphecyClaim)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CosmosBridgeLogNewProphecyClaim)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CosmosBridgeLogNewProphecyClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogNewProphecyClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogNewProphecyClaim represents a LogNewProphecyClaim event raised by the CosmosBridge contract.
type CosmosBridgeLogNewProphecyClaim struct {
	ProphecyID       *big.Int
	Nonce            *big.Int
	CosmosSender     []byte
	EthereumReceiver common.Address
	ValidatorAddress common.Address
	TokenAddress     common.Address
	Symbol           string
	Amount           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogNewProphecyClaim is a free log retrieval operation binding the contract event 0x48597fd2bac5835331b7653bb9b6284df26bb7fbd712b0059a42baf36e1b15e9.
//
// Solidity: event LogNewProphecyClaim(uint256 _prophecyID, uint256 _nonce, bytes _cosmosSender, address _ethereumReceiver, address _validatorAddress, address _tokenAddress, string _symbol, uint256 _amount)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogNewProphecyClaim(opts *bind.FilterOpts) (*CosmosBridgeLogNewProphecyClaimIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogNewProphecyClaim")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogNewProphecyClaimIterator{contract: _CosmosBridge.contract, event: "LogNewProphecyClaim", logs: logs, sub: sub}, nil
}

// WatchLogNewProphecyClaim is a free log subscription operation binding the contract event 0x48597fd2bac5835331b7653bb9b6284df26bb7fbd712b0059a42baf36e1b15e9.
//
// Solidity: event LogNewProphecyClaim(uint256 _prophecyID, uint256 _nonce, bytes _cosmosSender, address _ethereumReceiver, address _validatorAddress, address _tokenAddress, string _symbol, uint256 _amount)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogNewProphecyClaim(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogNewProphecyClaim) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogNewProphecyClaim")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogNewProphecyClaim)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogNewProphecyClaim", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// CosmosBridgeLogOracleSetIterator is returned from FilterLogOracleSet and is used to iterate over the raw logs and unpacked data for LogOracleSet events raised by the CosmosBridge contract.
type CosmosBridgeLogOracleSetIterator struct {
	Event *CosmosBridgeLogOracleSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CosmosBridgeLogOracleSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogOracleSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CosmosBridgeLogOracleSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CosmosBridgeLogOracleSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogOracleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogOracleSet represents a LogOracleSet event raised by the CosmosBridge contract.
type CosmosBridgeLogOracleSet struct {
	Oracle common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogOracleSet is a free log retrieval operation binding the contract event 0x6efb0434342713e2e9b1501dbebf76b4ed18406ea77ab5d56535cc26dec3adc0.
//
// Solidity: event LogOracleSet(address _oracle)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogOracleSet(opts *bind.FilterOpts) (*CosmosBridgeLogOracleSetIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogOracleSet")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogOracleSetIterator{contract: _CosmosBridge.contract, event: "LogOracleSet", logs: logs, sub: sub}, nil
}

// WatchLogOracleSet is a free log subscription operation binding the contract event 0x6efb0434342713e2e9b1501dbebf76b4ed18406ea77ab5d56535cc26dec3adc0.
//
// Solidity: event LogOracleSet(address _oracle)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogOracleSet(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogOracleSet) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogOracleSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogOracleSet)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogOracleSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// CosmosBridgeLogProphecyCompletedIterator is returned from FilterLogProphecyCompleted and is used to iterate over the raw logs and unpacked data for LogProphecyCompleted events raised by the CosmosBridge contract.
type CosmosBridgeLogProphecyCompletedIterator struct {
	Event *CosmosBridgeLogProphecyCompleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CosmosBridgeLogProphecyCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogProphecyCompleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CosmosBridgeLogProphecyCompleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CosmosBridgeLogProphecyCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogProphecyCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogProphecyCompleted represents a LogProphecyCompleted event raised by the CosmosBridge contract.
type CosmosBridgeLogProphecyCompleted struct {
	ProphecyID *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogProphecyCompleted is a free log retrieval operation binding the contract event 0x18083372d5d670784eb4cddbb916a6b39bae1015476bdfe19a0fb1482c45313d.
//
// Solidity: event LogProphecyCompleted(uint256 _prophecyID)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogProphecyCompleted(opts *bind.FilterOpts) (*CosmosBridgeLogProphecyCompletedIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogProphecyCompleted")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogProphecyCompletedIterator{contract: _CosmosBridge.contract, event: "LogProphecyCompleted", logs: logs, sub: sub}, nil
}

// WatchLogProphecyCompleted is a free log subscription operation binding the contract event 0x18083372d5d670784eb4cddbb916a6b39bae1015476bdfe19a0fb1482c45313d.
//
// Solidity: event LogProphecyCompleted(uint256 _prophecyID)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogProphecyCompleted(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogProphecyCompleted) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogProphecyCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogProphecyCompleted)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogProphecyCompleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
