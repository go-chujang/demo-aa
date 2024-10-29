package aa

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/common/logx"
	"github.com/go-chujang/demo-aa/common/utils/ternary"
	"github.com/go-chujang/demo-aa/config"
	"github.com/go-chujang/demo-aa/internal/storedquery"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx/rpcx"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

const (
	watchDogMinInterval   = time.Millisecond * 500
	watchDogMinBlockRange = uint64(1)
)

var defaultWatchAbis = []abi.ABI{
	EntrypointAbi,
	AccountFactoryAbi,
	PaymasterAbi,
	GambleAbi,
}

type WatchDog struct {
	db         *mongox.Client
	rpcUri     string
	interval   time.Duration
	blockRange uint64

	ctx    context.Context
	cancel context.CancelFunc

	addresses []common.Address
	EventMap  map[string]abi.Event
	ErrorMap  map[string]abi.Error
}

func NewWatchDog(db *mongox.Client, rpcUri string, interval time.Duration, blockRange uint64, abiOps ...abi.ABI) (*WatchDog, error) {
	mngd, err := storedquery.GetMngdContracts(db)
	if err != nil {
		return nil, err
	}
	watchDog := &WatchDog{
		db:         db,
		rpcUri:     rpcUri,
		interval:   ternary.Default(nil, interval, watchDogMinInterval),
		blockRange: ternary.Default(nil, blockRange, watchDogMinBlockRange),
		addresses:  make([]common.Address, 0, len(mngd)),
		EventMap:   make(map[string]abi.Event),
		ErrorMap:   make(map[string]abi.Error),
	}
	for _, v := range mngd {
		watchDog.addresses = append(watchDog.addresses, v.Address)
	}
	for _, a := range ternary.VArgs(nil, defaultWatchAbis, abiOps) {
		for _, ev := range a.Events {
			eventId := ev.ID.String()
			watchDog.EventMap[eventId] = ev
		}
		for _, er := range a.Errors {
			errId := er.ID.String()
			watchDog.ErrorMap[errId] = er
		}
	}
	return watchDog, nil
}

func (wd *WatchDog) IsStopped() bool {
	return wd.ctx == nil
}

func (wd *WatchDog) Stop() error {
	if !wd.IsStopped() {
		wd.cancel()
		wd.ctx = nil
		wd.cancel = nil
	}
	return nil
}

func (wd *WatchDog) Start() {
	wd.Stop()
	startingBlockNumber, err := storedquery.BlockNumberInc(wd.db, 0)
	if err != nil {
		panic(err)
	}
	logx.Debug(config.AppTag(), "startingBlockNumber: %d", startingBlockNumber)

	wd.ctx, wd.cancel = context.WithCancel(context.Background())
	defer func() {
		if r := recover(); r != nil {
			logx.Criticalf("recovered WatchDog: %v", r)
			time.Sleep(3 * time.Second)
			wd.Start()
		}
	}()
	wd.watch(wd.ctx, startingBlockNumber)
}

func (wd *WatchDog) watch(ctx context.Context, startingBlockNumber uint64) {
	ticker := time.NewTicker(time.Duration(wd.interval))
	defer ticker.Stop()

	var (
		noRace          atomic.Bool
		lastBlockNumber = startingBlockNumber
	)
	taskFn := func() error {
		if noRace.Load() {
			return nil
		}
		noRace.Store(true)
		defer noRace.Store(false)

		current, err := rpcx.EasyUint(wd.rpcUri, rpcx.MethodBlockNumber)
		if err != nil {
			return fmt.Errorf("watch::get current %s", err.Error())
		}

		if current <= lastBlockNumber {
			lastBlockNumber = current
			return nil
		}

		gap := current - lastBlockNumber
		if gap > wd.blockRange {
			gap = wd.blockRange
		}
		toBlock := lastBlockNumber + gap

		elogs, err := rpcx.FilterLogs2[model.EmittedLog](wd.rpcUri, rpcx.FilterQuery{
			FromBlock: new(big.Int).SetUint64(lastBlockNumber),
			ToBlock:   new(big.Int).SetUint64(toBlock),
			Addresses: wd.addresses,
		})
		if err != nil {
			return fmt.Errorf("watch::get emitted logs %s", err.Error())
		}
		if err := wd.writeToDB(elogs); err != nil {
			return fmt.Errorf("watch::write to db %s", err.Error())
		}
		if _, err := storedquery.BlockNumberInc(wd.db, gap); err != nil {
			return fmt.Errorf("watch::update blocknumber %s", err.Error())
		}
		lastBlockNumber = toBlock
		return nil
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := taskFn(); err != nil {
				logx.Errorf(err.Error())
			}
		}
	}
}

func (wd *WatchDog) writeToDB(queriedLogs []model.EmittedLog) error {
	if len(queriedLogs) == 0 {
		return nil
	}
	var (
		writeModels = make(map[string][]mongo.WriteModel)
		deduplHash  = make(map[string]bool)
		mu          sync.Mutex
		eg          = new(errgroup.Group)
	)

	collEmitLog := model.EmittedLog{}.Collection()
	writeModels[collEmitLog] = make([]mongo.WriteModel, 0, len(queriedLogs))
	for _, v := range queriedLogs {
		log := v
		eg.Go(func() error {
			if err := log.ParseByAbi(wd.EventMap, wd.ErrorMap); err != nil {
				return err
			}
			replace := mongo.NewReplaceOneModel().
				SetFilter(bson.M{"_id": log.ID()}).
				SetReplacement(log).
				SetUpsert(true)

			mu.Lock()
			writeModels[collEmitLog] = append(writeModels[collEmitLog], replace)
			deduplHash[log.TransactionHash] = true
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	txHashes := make([]string, 0, len(deduplHash))
	for k := range deduplHash {
		txHashes = append(txHashes, k)
	}
	txLogs, err := rpcx.GetTxReceipts2[model.TxLog](wd.rpcUri, txHashes)
	if err != nil {
		return err
	}

	collTxLog := model.TxLog{}.Collection()
	writeModels[collTxLog] = make([]mongo.WriteModel, 0, len(txLogs))
	for _, v := range txLogs {
		writeModels[collTxLog] = append(writeModels[collTxLog], mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": v.ID()}).
			SetUpdate(bson.M{"$set": bson.M{
				"status":            v.Status,
				"from":              v.From.Hex(),
				"to":                v.To.Hex(),
				"blockNumber":       v.BlockNumber,
				"gasUsed":           v.GasUsed,
				"effectiveGasPrice": v.EffectiveGasPrice,
			}}).
			SetUpsert(true),
		)
	}
	query := mongox.NewQuery().SetOrdered(false)
	return wd.db.BulkWriteMultiCollections(writeModels, query, wd.ctx)
}
