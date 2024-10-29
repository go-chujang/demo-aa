package storedquery

import (
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"go.mongodb.org/mongo-driver/bson"
)

func GetMngdWallet(db *mongox.Client) (*model.MngdWallet, error) {
	var (
		wallet   model.MngdWallet
		accounts []model.ManagedAccount
		query    = mongox.NewQuery().SetColl(model.CollectionMngdAccounts).SetFilter(bson.M{
			"kind": model.KindEOA,
		})
	)
	if err := db.Find(&accounts, query); err != nil {
		return nil, err
	}
	for _, v := range accounts {
		switch *v.Role {
		case model.RoleSupervisor:
			wallet.Supervisor = v
		case model.RoleManager:
			wallet.Manager = v
		case model.RoleOperator:
			wallet.Operators = append(wallet.Operators, v)
		}
	}
	return &wallet, nil
}

func GetMngdContracts(db *mongox.Client) ([]model.ManagedAccount, error) {
	var (
		mngd  []model.ManagedAccount
		query = mongox.NewQuery().SetColl(model.CollectionMngdAccounts).SetFilter(bson.M{
			"kind": model.KindContract,
		})
	)
	return mngd, db.Find(&mngd, query)
}
