package kafka

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/go-chujang/demo-aa/common/logx"
	"github.com/go-chujang/demo-aa/common/utils/slice"
)

func Produce(topic Topic, key, value sarama.Encoder) error {
	producerMustNotNil()

	producer, err := producerPool.borrow()
	if err != nil {
		return err
	}
	defer producerPool.release(producer)

	if err = producer.BeginTxn(); err != nil {
		return err
	}

	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic.String(),
		Key:   key,
		Value: value,
	}

	err = producer.CommitTxn()
	switch {
	case producer.TxnStatus()&sarama.ProducerTxnFlagFatalError != 0:
		err = errors.Join(err, errors.New("failed commit with txn fatal error"))
	case producer.TxnStatus()&sarama.ProducerTxnFlagAbortableError != 0:
		if abortErr := producer.AbortTxn(); abortErr != nil {
			err = errors.Join(err, abortErr)
		}
	}
	return err
}

func Consume(ctx context.Context, topicNames []Topic, batchHandler ConsumeBatchHandler) error {
	consumerMustNotNil()

	switch {
	case ctx == nil:
		return errors.New("context must not be nil")
	case batchHandler == nil:
		return errors.New("batchHandler must not be nil")
	default:
		if topicNames == nil {
			return errors.New("topics must not be nil")
		}
	}
	consumer.batchHandler = batchHandler
	topics := slice.TypeCast(topicNames, func(o Topic) string { return o.String() })

	for {
		if err := consumer.client.Consume(ctx, topics, consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return err
			}
			logx.Error(consumer.groupId, "topic name: %v, error: %s", topics, err.Error())
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		consumer.Ready = make(chan bool)
	}
}
