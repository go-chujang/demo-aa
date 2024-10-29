package mongox

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidOperation     = errors.New("invalid op")
	ErrEmptyDbCollection    = errors.New("empty db or collection")
	ErrMismatchedCollection = errors.New("mismatched collection")
	ErrNotFoundDbCollection = errors.New("not found db or collection")
	ErrNoMatchedDocuments   = errors.New("no documents matched the update filter")
	ErrZeroModifiedUpdate   = errors.New("zero modified document in update")
)

func updateError(err error, res *mongo.UpdateResult, isUpsert bool) error {
	if err != nil {
		return err
	}
	switch {
	case res.MatchedCount == 0:
		return ErrNoMatchedDocuments
	case isUpsert && res.UpsertedCount == 0 && res.ModifiedCount == 0:
	case !isUpsert && res.ModifiedCount == 0:
	default:
		return nil
	}
	return ErrZeroModifiedUpdate
}
