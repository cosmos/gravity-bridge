package utils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	TestNullAddress  = "0x0000000000000000000000000000000000000000"
	TestOtherAddress = "0x1000000000000000000000000000000000000000"
	TestCoinString   = "9921ruby"
	TestCoinSymbol   = "ruby"
	TestCoinAmount   = 9921
)

func TestIsZeroAddress(t *testing.T) {
	falseRes := IsZeroAddress(common.HexToAddress(TestOtherAddress))
	require.False(t, falseRes)

	trueRes := IsZeroAddress(common.HexToAddress(TestNullAddress))
	require.True(t, trueRes)
}

func TestGetSymbolFromCoin(t *testing.T) {
	testCoinAmount := big.NewInt(int64(TestCoinAmount))

	symbol, amount := GetSymbolAmountFromCoin(TestCoinString)

	require.Equal(t, TestCoinSymbol, symbol)
	require.Equal(t, testCoinAmount, amount)
}
