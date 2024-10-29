package id

import "go.mongodb.org/mongo-driver/bson/primitive"

func Hex() string {
	return primitive.NewObjectID().Hex()
}
