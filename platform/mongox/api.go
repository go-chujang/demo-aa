package mongox

import (
	"context"
	"errors"

	"github.com/go-chujang/demo-aa/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Client) collCtx(query *Query, ctxorigin ...context.Context) (*mongo.Collection, context.Context, error) {
	if query.Database == "" || query.Collection == "" {
		return nil, nil, ErrEmptyDbCollection
	}
	return c.client.Database(query.Database).Collection(query.Collection), c.ctx(ctxorigin...), nil
}

func (c *Client) Count(query *Query, ctxorigin ...context.Context) (int64, error) {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return 0, err
	}
	return coll.CountDocuments(ctx, query.Filter)
}

func (c *Client) Find(docs interface{}, query *Query, ctxorigin ...context.Context) error {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return err
	}
	filter, opts := query.find()
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, docs)
}

func (c *Client) FindOne(doc interface{}, query *Query, ctxorigin ...context.Context) error {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return err
	}
	filter, opts := query.findOne()
	return coll.FindOne(ctx, filter, opts).Decode(doc)
}

func (c *Client) Aggregate(docs interface{}, query *Query, ctxorigin ...context.Context) error {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return err
	}
	pipeline, opts := query.aggregate()
	cursor, err := coll.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, docs)
}

func (c *Client) Exist(query *Query, ctxorigin ...context.Context) (bool, error) {
	err := c.FindOne(&bson.M{}, query, ctxorigin...)
	switch err {
	case nil: // exist
		return true, nil
	case mongo.ErrNoDocuments: // not exist
		return false, nil
	default:
		return false, err
	}
}

func (c *Client) InsertOne(doc model.Document, query *Query, ctxorigin ...context.Context) error {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, doc)
	return err
}

func (c *Client) InsertMany(docs []model.Document, query *Query, ctxorigin ...context.Context) ([]interface{}, error) {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return nil, err
	}

	documents := make([]interface{}, 0, len(docs))
	for _, v := range docs {
		if coll.Name() != v.Collection() {
			return nil, ErrMismatchedCollection
		}
		documents = append(documents, v)
	}
	res, err := coll.InsertMany(ctx, documents, query.insertMany())
	return res.InsertedIDs, err
}

func (c *Client) FindOneAndUpdate(doc interface{}, query *Query, ctxorigin ...context.Context) error {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return err
	}
	filter, update, opts := query.findOneAndUpdate()
	return coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(doc)
}

func (c *Client) ReplaceOne(doc model.Document, query *Query, ctxorigin ...context.Context) error {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return err
	}
	filter, opts := query.replaceOne()
	res, err := coll.ReplaceOne(ctx, filter, doc, opts)
	return updateError(err, res, *opts.Upsert)
}

func (c *Client) UpdateOne(query *Query, ctxorigin ...context.Context) error {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return err
	}
	filter, update, opts := query.update()
	res, err := coll.UpdateOne(ctx, filter, update, opts)
	return updateError(err, res, *opts.Upsert)
}

func (c *Client) UpdateMany(query *Query, ctxorigin ...context.Context) error {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return err
	}
	filter, update, opts := query.update()
	res, err := coll.UpdateMany(ctx, filter, update, opts)
	return updateError(err, res, *opts.Upsert)
}

func (c *Client) BulkWrite(writeModels []mongo.WriteModel, query *Query, ctxorigin ...context.Context) ([]error, error) {
	coll, ctx, err := c.collCtx(query, ctxorigin...)
	if err != nil {
		return nil, err
	}
	opts := query.bulkWrite()
	errs := make([]error, len(writeModels))

	if _, err = coll.BulkWrite(ctx, writeModels, opts); err != nil {
		if bulkErr, ok := err.(mongo.BulkWriteException); ok {
			for _, we := range bulkErr.WriteErrors {
				errs[we.Index] = we.WriteError
			}
		}
	}
	return errs, err
}

// todo: errors per collection
func (c *Client) BulkWriteMultiCollections(writeModels map[string][]mongo.WriteModel, query *Query, ctxorigin ...context.Context) error {
	var (
		db    = c.client.Database(query.Database)
		ctx   = c.ctx(ctxorigin...)
		opts  = query.bulkWrite()
		bwErr []error
	)
	for name, writeModels := range writeModels {
		_, err := db.Collection(name).BulkWrite(ctx, writeModels, opts)
		if err != nil && *opts.Ordered {
			return err
		}
		bwErr = append(bwErr, err)
	}
	return errors.Join(bwErr...)
}
