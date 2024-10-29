package model

import (
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chujang/demo-aa/platform/etherx/ethutil"
	"github.com/go-chujang/demo-aa/platform/kafka"
)

var (
	_ Message       = (*CreateAccount)(nil)
	_ RawDataHelper = (*CreateAccount)(nil)
)

///////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////

type CreateAccount struct {
	UserId string         `json:"userId"`
	Owner  common.Address `json:"owner"`
}

func (m CreateAccount) Topic() kafka.Topic { return kafka.TopicMembership }
func (m CreateAccount) KeyValue() (key sarama.Encoder, value sarama.Encoder, err error) {
	if m.Owner.Cmp(ethutil.ZeroAddress) == 0 {
		return nil, nil, errors.New("address must not be zero")
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, nil, err
	}
	return sarama.ByteEncoder(m.Owner.Bytes()), sarama.ByteEncoder(bytes), nil
}

func (m *CreateAccount) Parse(msg *sarama.ConsumerMessage) error {
	return json.Unmarshal(msg.Value, &m)
}

func (m CreateAccount) Hint() *string { return toHintString(&m) }
func (m CreateAccount) RawData() (hint *string, data map[string]interface{}) {
	return m.Hint(), toHexBsonMap(map[string]interface{}{
		"userId": m.UserId,
		"owner":  m.Owner,
		"salt":   0,
	})
}
