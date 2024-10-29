package storedquery

import (
	"github.com/go-chujang/demo-aa/model"
	"github.com/go-chujang/demo-aa/platform/mongox"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUser(db *mongox.Client, userId string) (*model.UserAccount, error) {
	query := mongox.NewQuery().SetColl(model.CollectionUserAccounts).SetFilter(bson.M{"_id": userId})

	var user *model.UserAccount
	return user, db.FindOne(&user, query)
}

func PostProduce(db *mongox.Client, userId string) error {
	query := mongox.NewQuery().SetColl(model.CollectionUserAccounts).
		SetFilter(bson.M{"_id": userId}).
		SetUpdate(bson.M{"pending": true})
	return db.UpdateOne(query)
}
