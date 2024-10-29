package aa

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	hintCreateAccount = model.CreateAccount{}.Hint()
	hintUserOperation = model.PackedUserOperation{}.Hint()
)

func syncCreateAccount(before *model.UserAccount) (*model.UserAccount, error) {
	var after model.UserAccount
	var txlog model.TxLog

	queryTxLog := mongox.NewQuery().SetColl(model.CollectionTxLogs).SetFilter(bson.M{
		"blockNumber":    bson.M{"$gt": before.SyncedBlockNumber},
		"hint":           hintCreateAccount,
		"rawData.userId": before.UserId,
	})
	if err := api.db.FindOne(&txlog, queryTxLog); err != nil {
		return before, err
	}
	switch txlog.Status {
	case model.TxStatusSuccess:
	case model.TxStatusFailure: // todo: redrive
	default:
		return before, nil
	}

	owner, _ := txlog.RawData["owner"].(string)
	ownerAddr := common.HexToAddress(owner)
	account, err := api.getAddress(ownerAddr)
	if err != nil || ownerAddr == ethutil.ZeroAddress {
		return before, err
	}
	query := mongox.NewQuery().SetColl(model.CollectionUserAccounts).
		SetFilter(bson.M{"_id": before.ID()}).
		SetUpdate(bson.M{"$set": bson.M{
			"owner":             ownerAddr.Hex(),
			"account":           account.Hex(),
			"syncedBlockNumber": txlog.BlockNumber,
		}}).
		SetReturnAfter(true)
	if err = api.db.FindOneAndUpdate(&after, query); err != nil {
		return before, err
	}
	return &after, nil
}

func syncUserOperation(before *model.UserAccount) (*model.UserAccount, error) {
	var after model.UserAccount
	var txlog model.TxLog

	queryTxLog := mongox.NewQuery().SetColl(model.CollectionTxLogs).SetFilter(bson.M{
		"_id":  before.LastTxnHash,
		"hint": hintUserOperation,
	})
	if err := api.db.FindOne(&txlog, queryTxLog); err != nil {
		return before, err
	}
	switch txlog.Status {
	case model.TxStatusSuccess:
	case model.TxStatusFailure: // todo: redrive
	default:
		return before, nil
	}
	query := mongox.NewQuery().SetColl(model.CollectionUserAccounts).
		SetFilter(bson.M{"_id": before.ID()}).
		SetUpdate(bson.M{"$set": bson.M{
			"pending":           false,
			"syncedBlockNumber": txlog.BlockNumber,
		}}).
		SetReturnAfter(true)
	if err := api.db.FindOneAndUpdate(&after, query); err != nil {
		return before, err
	}
	return &after, nil
}
