package mongox

import (
	"sync"

	"github.com/go-chujang/demo-aa/config"
	"go.mongodb.org/mongo-driver/bson"
)

type Query struct {
	Database   string   `json:"database"`
	Collection string   `json:"collection"`
	Filter     bson.M   `json:"filter,omitempty"`
	UpdateSet  bson.M   `json:"updateSet,omitempty"`
	Pipeline   []bson.M `json:"pipeline,omitempty"`
	Options    struct {
		Sort        *bson.D `json:"sort,omitempty"`
		Limit       *int64  `json:"limit,omitempty"`
		Skip        *int64  `json:"skip,omitempty"`
		Upsert      bool    `json:"upsert,omitempty"`
		ReturnAfter bool    `json:"returnAfter,omitempty"`
		Project     bson.M  `json:"project,omitempty"`
		Hint        *string `json:"hint,omitempty"`
		Ordered     *bool   `json:"ordered,omitempty"`
	} `json:"options,omitempty"`
}

var (
	// enable empty string
	defaultDatabase     string
	defaultDatabaseOnce sync.Once
)

func NewQuery(collection ...string) *Query {
	defaultDatabaseOnce.Do(func() {
		defaultDatabase = config.Get(config.MONGO_DATABASE)
	})
	if collection != nil {
		return &Query{Database: defaultDatabase, Collection: collection[0]}
	}
	return &Query{Database: defaultDatabase}
}

func (q *Query) SetDB(db string) *Query {
	q.Database = db
	return q
}

func (q *Query) SetColl(collection string) *Query {
	q.Collection = collection
	return q
}

func (q *Query) SetFilter(filter bson.M) *Query {
	q.Filter = filter
	return q
}

func (q *Query) SetUpdate(updateSet bson.M) *Query {
	q.UpdateSet = updateSet
	return q
}

func (q *Query) SetPipeline(pipeline []bson.M) *Query {
	q.Pipeline = pipeline
	return q
}

func (q *Query) SetSort(sort bson.D) *Query {
	q.Options.Sort = &sort
	return q
}

func (q *Query) SetLimit(i int64) *Query {
	q.Options.Limit = &i
	return q
}

func (q *Query) SetSkip(i int64) *Query {
	q.Options.Skip = &i
	return q
}

func (q *Query) SetUpsert(upsert bool) *Query {
	q.Options.Upsert = upsert
	return q
}

func (q *Query) SetReturnAfter(after bool) *Query {
	q.Options.ReturnAfter = after
	return q
}

func (q *Query) SetProject(proj bson.M) *Query {
	q.Options.Project = proj
	return q
}

func (q *Query) SetHint(hint string) *Query {
	q.Options.Hint = &hint
	return q
}

func (q *Query) SetOrdered(ordered bool) *Query {
	q.Options.Ordered = &ordered
	return q
}
