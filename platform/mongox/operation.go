package mongox

type Operation struct {
	Op        string      `json:"op"`
	Query     Query       `json:"query"`
	Documents interface{} `json:"documents,omitempty"`
}

const (
	// op read
	op_count     = "count"
	op_find      = "find"
	op_findone   = "findone"
	op_aggregate = "aggregate"
	// op write
	op_insert     = "insert"
	op_insertone  = "insertone"
	op_replaceone = "replaceone"
	op_updateone  = "updateone"
)

func (o Operation) check() error {
	switch o.Op {
	// read
	case op_count:
	case op_find:
	case op_findone:
	case op_aggregate:
	// write
	case op_insert:
	case op_insertone:
	case op_replaceone:
	case op_updateone:
	default:
		return ErrInvalidOperation
	}
	return nil
}
