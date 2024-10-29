package storedquery

import (
	"github.com/go-chujang/demo-aa/common/utils/ternary"
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"go.mongodb.org/mongo-driver/bson"
)

func BlockNumberInc(db *mongox.Client, inc uint64, idOps ...string) (uint64, error) {
	param := model.Parameter{
		Id: ternary.VArgs(nil, "blockNumber", idOps...),
	}
	query := mongox.NewQuery().
		SetColl(param.Collection()).
		SetFilter(bson.M{"_id": param.Id}).
		SetUpdate(bson.M{"$inc": bson.M{"seq": inc}}).
		SetUpsert(true).
		SetReturnAfter(true)
	if err := db.FindOneAndUpdate(&param, query); err != nil {
		return 0, err
	}
	return *param.Seq, nil
}

// func FaucetAmount(db *mongox.Client, amountOps ...*big.Int) (*big.Int, error) {
// 	param := model.Parameter{Id: "faucetAmount"}
// 	if amountOps != nil {
// 		param.Parameter = amountOps[0]
// 	}
// 	query := mongox.NewQuery().
// 		SetColl(param.Collection()).
// 		SetFilter(bson.M{"_id": param.Id}).
// 		SetUpdate(bson.M{"$set": param}).
// 		SetUpsert(true).
// 		SetReturnAfter(true)
// 	if err := db.FindOneAndUpdate(&param, query); err != nil {
// 		return nil, err
// 	}
// 	var amount *big.Int
// 	return amount, param.V(&amount)
// }
