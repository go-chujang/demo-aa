package model

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/IBM/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/go-chujang/demo-aa/platform/kafka"
)

var (
	_ Message       = (*Faucet)(nil)
	_ RawDataHelper = (*Faucet)(nil)
)

type Faucet struct {
	Receiver common.Address `json:"receiver"`
	Value    *big.Int       `json:"value"`
}

func (m Faucet) Topic() kafka.Topic { return kafka.TopicLiquidity }
func (m Faucet) KeyValue() (key sarama.Encoder, value sarama.Encoder, err error) {
	if m.Receiver.Cmp(ethutil.ZeroAddress) == 0 {
		return nil, nil, errors.New("address must not be zero")
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, nil, err
	}
	return sarama.ByteEncoder(m.Receiver.Bytes()), sarama.ByteEncoder(bytes), nil
}

func (m *Faucet) Parse(msg *sarama.ConsumerMessage) error {
	return json.Unmarshal(msg.Value, &m)
}

func (m Faucet) Hint() *string { return toHintString(&m) }
func (m Faucet) RawData() (hint *string, data map[string]interface{}) {
	return m.Hint(), toHexBsonMap(map[string]interface{}{
		"userId": m.Receiver,
		"value":  m.Value,
	})
}
