package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	_ Document         = (*TxLog)(nil)
	_ json.Unmarshaler = (*TxLog)(nil)
	_ bson.Marshaler   = (*TxLog)(nil)
	_ bson.Unmarshaler = (*TxLog)(nil)
	_ json.Unmarshaler = (*txStatus)(nil)
)

type txStatus int

const (
	CollectionTxLogs          = "tx_logs"
	TxStatusPending  txStatus = 0
	TxStatusSuccess  txStatus = 1
	TxStatusFailure  txStatus = 2
)

// func (t txStatus) String() string { return string(t) }

type TxLog struct {
	Id                string         `bson:"_id" json:"transactionHash"` // transaction hash
	Status            txStatus       `bson:"status" json:"status"`
	From              common.Address `bson:"from" json:"from"`
	To                common.Address `bson:"to" json:"to"`
	BlockNumber       uint64         `bson:"blockNumber" json:"blockNumber"`
	GasUsed           uint64         `bson:"gasUsed" json:"gasUsed"`
	EffectiveGasPrice string         `bson:"effectiveGasPrice" json:"effectiveGasPrice"`

	Hint         *string                `bson:"hint,omitempty" json:"hint,omitempty"`
	RawData      map[string]interface{} `bson:"rawData,omitempty" json:"rawData,omitempty"`
	FailedReason *string                `bson:"failedReason,omitempty" json:"failedReason,omitempty"`
}

func (d TxLog) ID() string         { return d.Id }
func (d TxLog) Collection() string { return CollectionTxLogs }

// https://docs.infura.io/api/networks/ethereum/json-rpc-methods/eth_gettransactionreceipt
func (d *TxLog) UnmarshalJSON(input []byte) error {
	type TxLog struct {
		Id                *string                `json:"transactionHash"` // transaction hash
		Status            *txStatus              `json:"status"`
		From              *common.Address        `json:"from"`
		To                *common.Address        `json:"to"`
		BlockNumber       *hexutil.Uint64        `json:"blockNumber"`
		GasUsed           *hexutil.Uint64        `json:"gasUsed"`
		EffectiveGasPrice *hexutil.Big           `json:"effectiveGasPrice"`
		Hint              *string                `bson:"hint"`
		RawData           map[string]interface{} `json:"rawData"`
		FailedReason      *string                `json:"failedReason"`
	}
	var dec TxLog
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Id == nil {
		return errors.New("missing required field 'transactionHash' for TxLog")
	}
	d.Id = *dec.Id
	if dec.Status != nil {
		d.Status = *dec.Status
	}
	if dec.From == nil {
		return errors.New("missing required field 'from' for TxLog")
	}
	d.From = *dec.From
	if dec.To == nil {
		d.To = ethutil.ZeroAddress
	} else {
		d.To = *dec.To
	}
	if dec.BlockNumber != nil {
		d.BlockNumber = uint64(*dec.BlockNumber)
	}
	if dec.GasUsed != nil {
		d.GasUsed = uint64(*dec.GasUsed)
	}
	if dec.EffectiveGasPrice != nil {
		d.EffectiveGasPrice = dec.EffectiveGasPrice.String()
	}
	d.Hint = dec.Hint
	d.RawData = dec.RawData
	d.FailedReason = dec.FailedReason
	return nil
}

func (d TxLog) MarshalBSON() ([]byte, error) {
	if !ethutil.IsHexHash(d.Id) {
		return nil, ErrInvalidDocumentId
	}
	type TxLog struct {
		Id                string                 `bson:"_id"`
		Status            txStatus               `bson:"status"`
		From              string                 `bson:"from"`
		To                string                 `bson:"to"`
		BlockNumber       uint64                 `bson:"blockNumber"`
		GasUsed           uint64                 `bson:"gasUsed"`
		EffectiveGasPrice string                 `bson:"effectiveGasPrice"`
		Hint              *string                `bson:"hint,omitempty"`
		RawData           map[string]interface{} `bson:"rawData,omitempty"`
		FailedReason      *string                `bson:"failedReason,omitempty"`
	}
	return bson.Marshal(&TxLog{
		Id:                d.Id,
		Status:            d.Status,
		From:              d.From.Hex(),
		To:                d.To.Hex(),
		BlockNumber:       d.BlockNumber,
		GasUsed:           d.GasUsed,
		EffectiveGasPrice: d.EffectiveGasPrice,
		Hint:              d.Hint,
		RawData:           toHexBsonMap(d.RawData),
		FailedReason:      d.FailedReason,
	})
}

func (d *TxLog) UnmarshalBSON(data []byte) error {
	type TxLog struct {
		Id                string                 `bson:"_id"`
		Status            txStatus               `bson:"status"`
		From              string                 `bson:"from"`
		To                string                 `bson:"to"`
		BlockNumber       uint64                 `bson:"blockNumber"`
		GasUsed           uint64                 `bson:"gasUsed"`
		EffectiveGasPrice string                 `bson:"effectiveGasPrice"`
		Hint              *string                `bson:"hint,omitempty"`
		RawData           map[string]interface{} `bson:"rawData,omitempty"`
		FailedReason      *string                `bson:"failedReason,omitempty"`
	}
	var dec TxLog
	if err := bson.Unmarshal(data, &dec); err != nil {
		return err
	}
	d.Id = dec.Id
	d.Status = dec.Status
	d.From = common.HexToAddress(dec.From)
	d.To = common.HexToAddress(dec.To)
	d.BlockNumber = dec.BlockNumber
	d.GasUsed = dec.GasUsed
	d.EffectiveGasPrice = dec.EffectiveGasPrice
	d.Hint = dec.Hint
	d.RawData = dec.RawData
	d.FailedReason = dec.FailedReason
	return nil
}

func (t *txStatus) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err == nil {
		switch statusStr {
		case "0x0":
			*t = TxStatusFailure
		case "0x1":
			*t = TxStatusSuccess
		default:
			*t = TxStatusPending
		}
		return nil
	}
	var statusInt int
	if err := json.Unmarshal(data, &statusInt); err == nil {
		switch statusInt {
		case 0:
			*t = TxStatusPending
		case 1:
			*t = TxStatusSuccess
		case 2:
			*t = TxStatusFailure
		default:
			return fmt.Errorf("invalid 'status' value: %d", statusInt)
		}
		return nil
	}
	return errors.New("invalid 'status' field")
}
