package mongox

import (
	"github.com/go-chujang/demo-aa/common/utils/ternary"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (q Query) find() (bson.M, *options.FindOptions) {
	opts := options.Find()
	if q.Options.Limit != nil {
		opts.SetLimit(*q.Options.Limit)
	}
	if q.Options.Skip != nil {
		opts.SetSkip(*q.Options.Skip)
	}
	if q.Options.Sort != nil {
		opts.SetSort(*q.Options.Sort)
	}
	if q.Options.Project != nil {
		opts.SetProjection(q.Options.Project)
	}
	if q.Options.Hint != nil {
		opts.SetHint(*q.Options.Hint)
	}
	return q.Filter, opts
}

func (q Query) findOne() (bson.M, *options.FindOneOptions) {
	opts := options.FindOne()
	if q.Options.Sort != nil {
		opts.SetSort(*q.Options.Sort)
	}
	if q.Options.Project != nil {
		opts.SetProjection(q.Options.Project)
	}
	if q.Options.Hint != nil {
		opts.SetHint(*q.Options.Hint)
	}
	return q.Filter, opts
}

func (q Query) insertMany() *options.InsertManyOptions {
	opts := options.InsertMany()
	if q.Options.Ordered != nil {
		opts.SetOrdered(*q.Options.Ordered)
	}
	return opts
}

func (q Query) bulkWrite() *options.BulkWriteOptions {
	opts := options.BulkWrite()
	if q.Options.Ordered != nil {
		opts.SetOrdered(*q.Options.Ordered)
	}
	return opts
}

func (q Query) findOneAndUpdate() (bson.M, bson.M, *options.FindOneAndUpdateOptions) {
	opts := options.FindOneAndUpdate()
	if q.Options.Sort != nil {
		opts.SetSort(*q.Options.Sort)
	}
	rd := ternary.Cond(q.Options.ReturnAfter, options.After, options.Before)
	return q.Filter, q.UpdateSet, opts.
		SetUpsert(q.Options.Upsert).
		SetReturnDocument(rd)
}

func (q Query) update() (bson.M, bson.M, *options.UpdateOptions) {
	return q.Filter, q.UpdateSet, options.Update().SetUpsert(q.Options.Upsert)
}

func (q Query) replaceOne() (bson.M, *options.ReplaceOptions) {
	return q.Filter, options.Replace().SetUpsert(q.Options.Upsert)
}

func (q Query) aggregate() ([]bson.M, *options.AggregateOptions) {
	opts := options.Aggregate()
	return q.Pipeline, opts
}
