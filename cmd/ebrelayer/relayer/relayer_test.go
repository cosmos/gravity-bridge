package relayer

// ------------------------------------------------------------
//    Relayer_Test
//
//    Tests Relayer functionality.
//
// TODO: Object orient relayer so that these tests will work
// ------------------------------------------------------------

// import (
// 	"context"
// 	"errors"
// 	"math/big"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"

// )

// type testTrigger struct {
// 	shouldRun bool
// 	runErr    error
// }

// func (t *testTrigger) Description() string {
// 	return "testtrigger"
// }

// type lastBlockData struct {
// 	eventType       string
// 	contractAddress string
// 	lastBlockNumber uint64
// }

// type testSubscription struct {
// }

// func (t *testSubscription) Unsubscribe() {
// }

// func (t *testSubscription) Err() <-chan error {
// 	return make(chan error)
// }

// type testChainReader struct {
// }

// func (t *testChainReader) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
// 	return &types.Block{}, nil
// }
// func (t *testChainReader) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
// 	return &types.Block{}, nil
// }
// func (t *testChainReader) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
// 	return &types.Header{
// 		Time: big.NewInt(88888888),
// 	}, nil
// }
// func (t *testChainReader) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
// 	return &types.Header{
// 		Time: big.NewInt(88888888),
// 	}, nil
// }
// func (t *testChainReader) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
// 	return uint(0), nil
// }
// func (t *testChainReader) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
// 	return &types.Transaction{}, nil
// }
// func (t *testChainReader) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
// 	return &testSubscription{}, nil
// }

// func (n *testPersister) LastBlockNumber(eventType string, contractAddress common.Address) uint64 {
// 	return n.lastBlock.lastBlockNumber
// }

// func (n *testPersister) LastBlockHash(eventType string, contractAddress common.Address) common.Hash {
// 	return common.Hash{}
// }

// func (n *testPersister) UpdateLastBlockData(events []*model.Event) error {
// 	if len(events) == 0 {
// 		return n.updateLastBlockError
// 	}
// 	event := events[0]
// 	n.lastBlock.eventType = event.EventType()
// 	n.lastBlock.contractAddress = event.ContractAddress().Hex()
// 	rawLog := event.LogPayload()
// 	n.lastBlock.lastBlockNumber = rawLog.BlockNumber
// 	return n.updateLastBlockError
// }

// func (t *testErrorWatcher) ContractAddress() common.Address {
// 	return common.HexToAddress("")
// }

// func relayerStart(t *testing.T, errChan chan error) {
// 	err := initRelayer()
// 	if err != nil {
// 		t.Errorf("Error initializing relayer: err: %v", err)
// 		errChan <- err
// 	}
// }

// func setupTestRelayer(contracts contractData) relayerMock {

// 	testRelayer := relayer.Relayer(
// 		&relayer.Config{
// 			Chain:              &testChainReader{},
// 			WsClient:           contracts.Client,
// 			Triggers:           triggers,
// 			StartBlock:         0,
// 		},
// 	)
// 	return testRelayer
// }

// func TestNewRelayer(t *testing.T) {
// 	contracts, err := cutils.SetupAllTestContracts()
// 	if err != nil {
// 		t.Fatalf("Unable to setup the contracts: %v", err)
// 	}
// 	collector := setupTestRelayer(contracts)

// 	errChan := make(chan error)

// 	select {
// 	case err := <-errChan:
// 		t.Errorf("Should not have received error on start relaying: err: %v", err)
// 	case <-time.After(5 * time.Second):
// 	}
// }

// func TestEventRelay(t *testing.T) {
// 	contracts, err := cutils.SetupAllTestContracts()
// 	if err != nil {
// 		t.Fatalf("Unable to setup the contracts: %v", err)
// 	}
// 	collector, persister := setupTestCollectorTestPersister(contracts)

// 	errChan := make(chan error)
// 	go collectionStart(collector, t, errChan)

// 	<-collector.StartChan()
// 	_, err = contracts.PeggyContract.Apply(contracts.Auth, contracts.PeggyContract, big.NewInt(400), "")
// 	if err != nil {
// 		t.Fatalf("Application failed: err: %v", err)
// 	}

// 	contracts.Client.Commit()

// 	_, err = contracts.PeggyContract.Withdraw(contracts.Auth, contracts.PeggyContract, big.NewInt(50))
// 	if err != nil {
// 		t.Fatalf("Withdrawal failed: err: %v", err)
// 	}

// 	contracts.Client.Commit()

// 	_, err = contracts.PeggyContract.Send(contracts.Auth, contracts.PeggyContract, big.NewInt(50))
// 	if err != nil {
// 		t.Fatalf("Deposit failed: err: %v", err)
// 	}

// 	contracts.Client.Commit()

// 	// Sleep for a bit to make sure all the events gets handled and stored
// 	time.Sleep(4 * time.Second)

// 	events, _ := persister.RetrieveEvents(&model.RetrieveEventsCriteria{
// 		Offset:  0,
// 		Count:   10,
// 		Reverse: false,
// 	})

// 	if len(events) == 0 {
// 		t.Error("Should have seen some events in the persister")
// 	}

// 	if len(events) != 6 {
// 		t.Errorf("Should have seen 6 events in the persister, saw %v instead", len(events))
// 		for _, event := range events {
// 			t.Logf("event = %v", event.EventType())
// 		}
// 	}

// 	err = collector.StopCollection(true)
// 	if err != nil {
// 		t.Errorf("Should not have returned an error when stopping collection: err: %v", err)
// 	}
// }

// func TestCheckRetrievedEvents(t *testing.T) {
// 	contracts, err := cutils.SetupAllTestContracts()
// 	if err != nil {
// 		t.Fatalf("Unable to setup the contracts: %v", err)
// 	}
// 	collector, _ := setupTestCollectorTestPersisterBadUpdateBlockData(contracts)

// 	testAddress := "0xdfe273082089bb7f70ee36eebcde64832fe97e55"
// 	testApplicationWhitelisted := &contract.PeggyContract{
// 		ListingAddress: common.HexToAddress(testAddress),
// 		Raw: types.Log{
// 			Address:     common.HexToAddress(testAddress),
// 			Topics:      []common.Hash{},
// 			Data:        []byte{},
// 			BlockNumber: 8888888,
// 			Index:       1,
// 		},
// 	}

// 	_, err = relayer.CheckRetrievedEvents(pastEvents)
// 	if err != nil {
// 		t.Errorf("Error checking retrieved events: %v", err)
// 	}

// }
