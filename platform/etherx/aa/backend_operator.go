package aa

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx"
	"github.com/go-chujang/demo-aa/platform/etherx/rpcx"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	OperatorOccupyTimeoutSeconds int64  = 10
	handleOpsGasLimit            uint64 = 5_000_000
	bunleChunkSize               int    = 10
)

type operatorBackend struct {
	*accountAbstract

	occupied bool
	chainId  *big.Int
	txrId    string
	txrs     map[string]*etherx.Transactor
	mu       sync.Mutex
}

func NewOpBackend(db *mongox.Client, rpcUri string) (*operatorBackend, error) {
	aa, err := newAA(db, rpcUri)
	if err != nil {
		return nil, err
	}
	aa.entrypoint.SetMaxGasLimit(handleOpsGasLimit)

	chainId, err := rpcx.EasyBig(rpcUri, rpcx.MethodChainId)
	if err != nil {
		return nil, err
	}

	opTxr := &operatorBackend{}
	opTxr.accountAbstract = aa
	opTxr.chainId = chainId
	opTxr.txrs = make(map[string]*etherx.Transactor)
	for _, v := range aa.mngdWallet.Operators {
		if *v.Role == model.RoleOperator {
			txr, _ := etherx.NewTransactor(rpcUri, v.PrivateKey, chainId)
			opTxr.txrs[v.ID()] = txr
		}
	}
	return opTxr, nil
}

func (o *operatorBackend) occupy() error {
	if o.occupied {
		return errors.New("already occupied")
	}

	operator := &model.ManagedAccount{}
	now := time.Now().Unix()
	query := mongox.NewQuery().SetColl(operator.Collection()).
		SetFilter(bson.M{
			"role": model.RoleOperator,
			"$or": bson.A{
				bson.M{"occupied": false},
				bson.M{"occupiedAt": bson.M{"$lt": now - OperatorOccupyTimeoutSeconds}},
			}}).
		SetSort(bson.D{bson.E{Key: "occupiedAt", Value: 1}}).
		SetUpdate(bson.M{
			"$set": bson.M{
				"occupied":   true,
				"occupiedAt": now,
			},
		})
	if err := o.accountAbstract.db.FindOneAndUpdate(&operator, query); err != nil {
		return err
	}

	id := operator.ID()
	if _, exist := o.txrs[id]; !exist {
		txr, _ := etherx.NewTransactor(o.rpcUri, operator.PrivateKey, o.chainId)
		o.txrs[id] = txr
	}
	o.txrId = id
	o.occupied = true
	return nil
}

func (o *operatorBackend) getOccupiedTxr() (*etherx.Transactor, error) {
	if o.txrId == "" || !o.occupied {
		return nil, errors.New("failed getOccupiedTxr")
	}
	return o.txrs[o.txrId], nil
}

func (o *operatorBackend) release() {
	if !o.occupied {
		return
	}
	o.accountAbstract.db.UpdateOne(mongox.NewQuery(model.CollectionMngdAccounts).SetUpdate(bson.M{
		"$set": bson.M{"occupied": false},
	}))
	o.occupied = false
}

func (o *operatorBackend) prepareHandleOps(msg ...PackedUserOperation) []BundlePayload {
	total := len(msg)
	chunk := (total + bunleChunkSize - 1) / bunleChunkSize
	bundlePayloads := make([]BundlePayload, 0, chunk)

	for i := 0; i < chunk; i++ {
		start := i * bunleChunkSize
		end := start + bunleChunkSize
		if end > total {
			end = total
		}

		chunk := msg[start:end]
		packed, _ := o.packHandleOps(o.mngdWallet.Supervisor.Address, chunk...)
		bundlePayloads = append(bundlePayloads, BundlePayload{
			ca:     o.entrypoint.ContractAddress,
			packed: packed,
		})
	}
	return bundlePayloads
}

func (o *operatorBackend) BundleExec(ctx context.Context, userOpHint string, userOps ...PackedUserOperation) ([]error, error) {
	if userOps == nil {
		return nil, nil
	}
	payloads := o.prepareHandleOps(userOps...)

	o.mu.Lock()
	defer o.mu.Unlock()

	if err := o.occupy(); err != nil {
		return nil, err
	}
	defer o.release()

	txr, err := o.getOccupiedTxr()
	if err != nil {
		return nil, err
	}

	hashes, rpcErrs, err := o.bundleExec(txr, payloads)
	if err != nil {
		return nil, err
	}

	var (
		txUpdateSets   = make([]mongo.WriteModel, len(hashes))
		userOpSize     = len(userOps)
		userUpdateSets = make([]mongo.WriteModel, userOpSize)
	)
	for i, hash := range hashes {
		updateSet := bson.M{"hint": userOpHint}
		rpcerr := rpcErrs[i]
		if rpcErrs[i] != nil {
			updateSet["failedReason"] = rpcerr.Error()
		}
		txUpdateSets[i] = mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": hash}).
			SetUpdate(bson.M{"$set": updateSet}).
			SetUpsert(true)

		start := i * bunleChunkSize
		end := start + bunleChunkSize
		if end > userOpSize {
			end = userOpSize
		}

		for j := start; j < end; j++ {
			uop := userOps[j]
			updateSet := bson.M{
				"lastTxnHash": hash,
				"pending":     true,
			}
			if rpcerr == nil {
				updateSet["lastUsedNonce"] = uop.Nonce.Uint64()
			}
			userUpdateSets[j] = mongo.NewUpdateOneModel().
				SetFilter(bson.M{"account": uop.Sender.Hex()}).
				SetUpdate(bson.M{"$set": updateSet}).
				SetUpsert(false)
		}
	}
	txQuery := mongox.NewQuery().SetColl(model.CollectionTxLogs).SetOrdered(false)
	errs, err := o.db.BulkWrite(txUpdateSets, txQuery, ctx)
	if err != nil {
		return nil, err
	}

	userQuery := mongox.NewQuery().SetColl(model.CollectionUserAccounts).SetOrdered(false)
	userErrs, err := o.db.BulkWrite(userUpdateSets, userQuery, ctx)
	if err != nil {
		return errs, err
	}
	for i, v := range userErrs {
		e := errs[i]
		errs[i] = errors.Join(e, v)
	}
	return errs, err
}
