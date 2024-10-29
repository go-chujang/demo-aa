package model

import (
	"errors"

	"github.com/IBM/sarama"
	"github.com/go-chujang/demo-aa/platform/kafka"
)

// todo: PreInsert -> bson/marshaler & unmarshaler
type Document interface {
	ID() string
	Collection() string
}

// errors used by Document implements
var (
	ErrEmptyDocumentId   = errors.New("empty document id")
	ErrEmptyDocumentBody = errors.New("empty document body")
	ErrInvalidDocumentId = errors.New("invalid document id")
	ErrInsufficientField = errors.New("insufficient required field")
)

type Message interface {
	Topic() kafka.Topic
	KeyValue() (key sarama.Encoder, value sarama.Encoder, err error)
	Parse(msg *sarama.ConsumerMessage) error
}

type RawDataHelper interface {
	Hint() *string
	RawData() (hint *string, data map[string]interface{}) // hint can nil
}
