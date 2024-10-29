package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-chujang/demo-aa/common/utils/check"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"go.mongodb.org/mongo-driver/bson"
)

const CollectionEmitLogs = "emitted_logs"

var (
	_ Document         = (*EmittedLog)(nil)
	_ json.Unmarshaler = (*EmittedLog)(nil)
	_ bson.Marshaler   = (*EmittedLog)(nil)
	_ bson.Unmarshaler = (*EmittedLog)(nil)
)

type EmittedLog struct {
	Id               string `bson:"_id" json:"_id"`
	BlockNumber      uint64 `bson:"-" json:"blockNumber,omitempty"`
	TransactionIndex uint   `bson:"-" json:"transactionIndex,omitempty"`
	LogIndex         uint   `bson:"-" json:"logIndex,omitempty"`

	Removed         bool           `bson:"removed" json:"removed"`
	TransactionHash string         `bson:"transactionHash" json:"transactionHash"`
	Address         common.Address `bson:"address" json:"address"`
	Topics          []common.Hash  `bson:"-" json:"topics"`
	Data            []byte         `bson:"-" json:"data"`

	EventName   string                 `bson:"eventName" json:"eventName"`
	Parameters  map[string]interface{} `bson:"parameters,omitempty" json:"parameters,omitempty"`
	ErrorName   *string                `bson:"errorName,omitempty" json:"errorName,omitempty"`
	Errors      map[string]interface{} `bson:"errors,omitempty" json:"errors,omitempty"`
	Additionals map[string]interface{} `bson:"additionals,omitempty" json:"additionals,omitempty"`
}

func (d *EmittedLog) ID() string {
	if d.Id == "" {
		d.Id = fmt.Sprintf("%d#%d#%d", d.BlockNumber, d.TransactionIndex, d.LogIndex)
	}
	return d.Id
}
func (d EmittedLog) Collection() string { return CollectionEmitLogs }

func (d *EmittedLog) UnmarshalJSON(input []byte) error {
	type EmittedLog struct {
		BlockNumber *hexutil.Uint64 `json:"blockNumber"`
		TxIndex     *hexutil.Uint   `json:"transactionIndex"`
		Index       *hexutil.Uint   `json:"logIndex"`
		Removed     *bool           `json:"removed"`
		TxHash      *common.Hash    `json:"transactionHash"`
		Address     *common.Address `json:"address"`
		Topics      []common.Hash   `json:"topics"`
		Data        *hexutil.Bytes  `json:"data"`
	}
	var dec EmittedLog
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.BlockNumber == nil || dec.TxIndex == nil || dec.Index == nil {
		return errors.New("missing required field for EmittedLog Id")
	}
	d.BlockNumber = uint64(*dec.BlockNumber)
	d.TransactionIndex = uint(*dec.TxIndex)
	d.LogIndex = uint(*dec.Index)

	if dec.TxHash == nil {
		return errors.New("missing required field 'transactionHash' for EmittedLog")
	}
	d.TransactionHash = dec.TxHash.Hex()

	if dec.Removed != nil {
		d.Removed = *dec.Removed
	}
	if dec.Address == nil {
		return errors.New("missing required field 'address' for EmittedLog")
	}
	d.Address = *dec.Address
	d.Topics = dec.Topics
	if dec.Data == nil {
		return errors.New("missing required field 'data' for EmittedLog")
	}
	d.Data = *dec.Data
	return nil
}

func (d EmittedLog) MarshalBSON() ([]byte, error) {
	type EmittedLog struct {
		Id              string                 `bson:"_id"`
		Removed         bool                   `bson:"removed"`
		TransactionHash string                 `bson:"transactionHash"`
		Address         string                 `bson:"address"`
		EventName       string                 `bson:"eventName"`
		Parameters      map[string]interface{} `bson:"parameters,omitempty"`
		ErrorName       *string                `bson:"errorName,omitempty"`
		Errors          map[string]interface{} `bson:"errors,omitempty"`
		Additionals     map[string]interface{} `bson:"additionals,omitempty"`
	}
	return bson.Marshal(&EmittedLog{
		Id:              d.ID(),
		Removed:         d.Removed,
		TransactionHash: d.TransactionHash,
		Address:         d.Address.Hex(),
		EventName:       d.EventName,
		Parameters:      toHexBsonMap(d.Parameters),
		ErrorName:       d.ErrorName,
		Errors:          toHexBsonMap(d.Errors),
		Additionals:     toHexBsonMap(d.Additionals),
	})
}

func (l *EmittedLog) UnmarshalBSON(data []byte) error {
	type EmittedLog struct {
		Id              string                 `bson:"_id"`
		Removed         bool                   `bson:"removed"`
		TransactionHash string                 `bson:"transactionHash"`
		Address         string                 `bson:"address"`
		EventName       string                 `bson:"eventName"`
		Parameters      map[string]interface{} `bson:"parameters,omitempty"`
		ErrorName       *string                `bson:"errorName,omitempty"`
		Errors          map[string]interface{} `bson:"errors,omitempty"`
		Additionals     map[string]interface{} `bson:"additionals,omitempty"`
	}
	var dec EmittedLog
	if err := bson.Unmarshal(data, &dec); err != nil {
		return err
	}
	l.Id = dec.Id
	success := false
	for range 1 {
		split := strings.Split(dec.Id, "#")
		if len(split) != 3 {
			break
		}
		blockNumber, err := strconv.ParseUint(split[0], 10, 64)
		if err != nil {
			break
		}
		txIndex, err := strconv.ParseUint(split[1], 10, 32)
		if err != nil {
			break
		}
		logIndex, err := strconv.ParseUint(split[2], 10, 32)
		if err != nil {
			break
		}
		l.BlockNumber = blockNumber
		l.TransactionIndex = uint(txIndex)
		l.LogIndex = uint(logIndex)
		success = true
	}
	if !success {
		return ErrInvalidDocumentId
	}
	l.Removed = dec.Removed
	l.TransactionHash = dec.TransactionHash
	l.Address = common.HexToAddress(dec.Address)
	l.EventName = dec.EventName
	l.Parameters = dec.Parameters
	l.ErrorName = dec.ErrorName
	l.Errors = dec.Errors
	l.Additionals = dec.Additionals
	return nil
}

func (d *EmittedLog) ParseByAbi(eventMap map[string]abi.Event, errorMap map[string]abi.Error) error {
	var err error
	d.EventName, d.Parameters, err = ethutil.ParseLog(types.Log{
		Address: d.Address,
		Topics:  d.Topics,
		Data:    d.Data,
	}, eventMap)
	if err == nil {
		return nil
	}
	if !check.ErrorOr(err, ethutil.ErrParseLogNilTopics, ethutil.ErrParseLogUnknownEvent) {
		return err
	}

	var errName string
	errName, d.Errors, err = ethutil.ParseError(d.Data, errorMap)
	if err == nil {
		d.ErrorName = &errName
	}
	if check.ErrorOr(err, ethutil.ErrParseErrInvalidData, ethutil.ErrParseErrUnknownError) {
		err = nil
	}
	return err
}
