package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/app"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.Gravity

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false
	gravityApp := app.Setup(checkTx)

	suite.ctx = gravityApp.BaseApp.NewContext(checkTx, tmproto.Header{Height: 1})
	suite.app = gravityApp

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, gravityApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, gravityApp.GravityKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestBridgeIDCRUD() {
	id := []byte("id")
	suite.app.GravityKeeper.SetBridgeID(suite.ctx, id)
	returnedID := suite.app.GravityKeeper.GetBridgeID(suite.ctx)
	suite.Require().Equal(id, returnedID)
}

func (suite *KeeperTestSuite) TestEthAddressCRUD() {
	cosmosAddr, err := types.GenerateTestCosmosAddress()
	suite.Require().NoError(err)
	valAddr, err := sdk.ValAddressFromHex(cosmosAddr.String())
	suite.Require().NoError(err)
	ethAddr, err := types.GenerateTestEthAddress()
	suite.Require().NoError(err)

	suite.app.GravityKeeper.SetEthAddress(suite.ctx, valAddr, *ethAddr)
	returnedEthAddr := suite.app.GravityKeeper.GetEthAddress(suite.ctx, valAddr)
	suite.Require().Equal(ethAddr, returnedEthAddr, "didn't receive set ethereum address")
}

func (suite *KeeperTestSuite) TestOrchestratorValidatorCRUD() {
	cosmosAddr, err := types.GenerateTestCosmosAddress()
	suite.Require().NoError(err)
	valAddr, err := sdk.ValAddressFromHex(cosmosAddr.String())
	suite.Require().NoError(err)
	cosmosAddr2, err := types.GenerateTestCosmosAddress()
	suite.Require().NoError(err)
	orchAddr, err := sdk.AccAddressFromHex(cosmosAddr2.String())
	suite.Require().NoError(err)

	suite.app.GravityKeeper.SetOrchestratorValidator(suite.ctx, valAddr, orchAddr)
	returnedEthAddr := suite.app.GravityKeeper.GetOrchestratorValidator(suite.ctx, orchAddr)
	suite.Require().Equal(valAddr, returnedEthAddr, "didn't receive set validator address")
}

func (suite *KeeperTestSuite) TestEthereumInfoCRUD() {
	ethInfo := types.EthereumInfo{Timestamp: time.Now(), Height: 10}

	suite.app.GravityKeeper.SetEthereumInfo(suite.ctx, ethInfo)
	returnedEthInfo, ok := suite.app.GravityKeeper.GetEthereumInfo(suite.ctx)
	suite.Require().True(ok, "no ethereum info located")
	suite.Require().Equal(ethInfo, returnedEthInfo, "didn't receive set ethereum info")
}

func (suite *KeeperTestSuite) TestLastObservedEventNonceCRUD() {
	nonce := uint64(13)

	suite.app.GravityKeeper.SetLastObservedEventNonce(suite.ctx, nonce)
	returnedEventNonce := suite.app.GravityKeeper.GetLastObservedEventNonce(suite.ctx)
	suite.Require().Equal(nonce, returnedEventNonce, "didn't receive set event nonce")
}

func (suite *KeeperTestSuite) TestTransferTxCRUD() {
	amount := sdk.Coin{Denom: "testdenom", Amount: sdk.NewIntFromUint64(100)}
	fee := sdk.Coin{Denom: "testdenom", Amount: sdk.NewIntFromUint64(12)}
	tx := types.TransferTx{Nonce: 13, Sender: "sender", EthereumRecipient: "recipient", Erc20Token: amount, Erc20Fee: fee}

	txid := suite.app.GravityKeeper.SetTransferTx(suite.ctx, tx)
	returnedTx, ok := suite.app.GravityKeeper.GetTransferTx(suite.ctx, txid)
	suite.Require().True(ok, "no transfer tx located")
	suite.Require().Equal(tx, returnedTx, "didn't receive set transfer transaction")

	suite.app.GravityKeeper.DeleteTransferTx(suite.ctx, txid)
	_, ok = suite.app.GravityKeeper.GetTransferTx(suite.ctx, txid)
	suite.Require().True(ok, "deleted transfer tx was returned")
}

func (suite *KeeperTestSuite) TestGetTransferTxs() {
	amount := sdk.Coin{Denom: "testdenom", Amount: sdk.NewIntFromUint64(100)}
	fee := sdk.Coin{Denom: "testdenom", Amount: sdk.NewIntFromUint64(12)}
	tx0 := types.TransferTx{Nonce: 13, Sender: "sender1", EthereumRecipient: "recipient1", Erc20Token: amount, Erc20Fee: fee}
	tx1 := types.TransferTx{Nonce: 13, Sender: "sender2", EthereumRecipient: "recipient2", Erc20Token: amount, Erc20Fee: fee}
	tx2 := types.TransferTx{Nonce: 13, Sender: "sender3", EthereumRecipient: "recipient3", Erc20Token: amount, Erc20Fee: fee}

	suite.app.GravityKeeper.SetTransferTx(suite.ctx, tx0)
	suite.app.GravityKeeper.SetTransferTx(suite.ctx, tx1)
	suite.app.GravityKeeper.SetTransferTx(suite.ctx, tx2)

	txs := suite.app.GravityKeeper.GetTransferTxs(suite.ctx)
	suite.Require().Len(txs, 3, "incorrect number of responses")
	suite.Require().Equal(tx0, txs[0], "incorrect transaction returned")
	suite.Require().Equal(tx1, txs[1], "incorrect transaction returned")
	suite.Require().Equal(tx2, txs[2], "incorrect transaction returned")
}

func (suite *KeeperTestSuite) TestEthereumEventCRUD() {
	id := []byte("testid")
	event := types.DepositEvent{Nonce: 20, TokenContract: "contract", Amount: sdk.NewInt(30), EthereumSender: "sender", CosmosReceiver: "receiver", EthereumHeight: 40}

	suite.app.GravityKeeper.SetEthereumEvent(suite.ctx, id, &event)
	returnedEvent, ok := suite.app.GravityKeeper.GetEthereumEvent(suite.ctx, id)
	suite.Require().True(ok, "set ethereum event was not returned")
	suite.Require().Equal(event, returnedEvent)
}

func (suite *KeeperTestSuite) TestConfirmCRUD() {
	id := []byte("testid")
	confirm := types.ConfirmLogicCall{InvalidationID: []byte("invitationID"), InvalidationNonce: 10, EthSigner: "signer", Signature: []byte("signature")}

	suite.app.GravityKeeper.SetConfirm(suite.ctx, id, &confirm)
	returnedConfirm, ok := suite.app.GravityKeeper.GetConfirm(suite.ctx, id)
	suite.Require().True(ok, "set ethereum confirm was not returned")
	suite.Require().Equal(confirm, returnedConfirm)
}
