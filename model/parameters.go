package model

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/common/utils/conv"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionParameters = "parameters"
	onlyForSeqIncrement  = "_only_for_seq_increment"
)

var (
	_             Document         = (*Parameter)(nil)
	_             bson.Marshaler   = (*Parameter)(nil)
	_             bson.Unmarshaler = (*Parameter)(nil)
	pmKindString                   = conv.TypeOf("string")
	pmKindBool                     = conv.TypeOf(true)
	pmKindUint64                   = conv.TypeOf(uint64(0))
	pmKindInt64                    = conv.TypeOf(int64(0))
	pmKindBigint                   = conv.TypeOf(big.NewInt(0))
	pmKindAddress                  = conv.TypeOf(common.Address{})
)

type Parameter struct {
	Id        string      `bson:"_id" json:"_id"`
	Kind      string      `bson:"kind" json:"-"`
	Seq       *uint64     `bson:"seq,omitempty" json:"seq,omitempty"`
	Value     string      `bson:"value" json:"-"`
	Parameter interface{} `bson:"-" json:"value"`
}

func (d Parameter) ID() string         { return d.Id }
func (p Parameter) Collection() string { return CollectionParameters }
func (p *Parameter) IncSeq(inc uint64) *Parameter {
	atomic.AddUint64(p.Seq, inc)
	return p
}

func (p *Parameter) V(v interface{}) error {
	if v == nil || reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("v must be pointer or nil interface")
	}
	if reflect.ValueOf(v).Elem().Type() != reflect.ValueOf(p.Parameter).Type() {
		return errors.New("mismatched type v")
	}
	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(p.Parameter))
	return nil
}

func (d Parameter) MarshalBSON() ([]byte, error) {
	if d.Id == "" {
		return nil, ErrEmptyDocumentId
	}

	kind := conv.TypeOf(d.Parameter)
	switch v := d.Parameter.(type) {
	case string:
		d.Value = v
	case bool:
		d.Value = strconv.FormatBool(v)
	case uint64:
		d.Value = strconv.FormatUint(v, 10)
	case int64:
		d.Value = strconv.FormatInt(v, 10)
	case *big.Int:
		if v == nil {
			v = big.NewInt(0)
		}
		d.Value = v.String()
	case common.Address:
		d.Value = v.Hex()
	default:
		if d.Seq == nil {
			return nil, fmt.Errorf("unsupported parameter type")
		}
		d.Value = onlyForSeqIncrement
		kind = pmKindUint64
	}

	type Parameter struct {
		Id    string  `bson:"_id"`
		Kind  string  `bson:"kind"`
		Seq   *uint64 `bson:"seq,omitempty"`
		Value string  `bson:"value"`
	}
	return bson.Marshal(&Parameter{
		Id:    d.Id,
		Kind:  kind,
		Seq:   d.Seq,
		Value: d.Value,
	})
}

func (d *Parameter) UnmarshalBSON(data []byte) error {
	type Parameter struct {
		Id    string  `bson:"_id"`
		Kind  string  `bson:"kind"`
		Seq   *uint64 `bson:"seq,omitempty"`
		Value string  `bson:"value"`
	}
	var dec Parameter
	err := bson.Unmarshal(data, &dec)
	if err != nil {
		return err
	}
	d.Id = dec.Id
	d.Kind = dec.Kind
	d.Seq = dec.Seq
	d.Value = dec.Value

	if dec.Seq != nil {
		d.Kind = pmKindUint64
		d.Value = onlyForSeqIncrement
		return nil
	}

	var parsed interface{}
	switch d.Kind {
	case pmKindString:
		parsed = d.Value
	case pmKindBool:
		parsed, err = strconv.ParseBool(d.Value)
	case pmKindUint64:
		parsed, err = strconv.ParseUint(d.Value, 10, 64)
	case pmKindInt64:
		parsed, err = strconv.ParseInt(d.Value, 10, 64)
	case pmKindBigint:
		if b, ok := new(big.Int).SetString(d.Value, 10); !ok {
			return errors.New("failed to convert Value to *big.Int")
		} else {
			parsed = b
		}
	case pmKindAddress:
		if !ethutil.IsAddress(d.Value) {
			return errors.New("invalid address format")
		} else {
			parsed = common.HexToAddress(d.Value)
		}
	default:
		return fmt.Errorf("unsupported parameter type")
	}
	if err != nil {
		return err
	}
	d.Parameter = parsed
	return nil
}
