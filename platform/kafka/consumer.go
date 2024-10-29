package kafka

import (
	"errors"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-chujang/demo-aa/common/logx"
)

var _ sarama.ConsumerGroupHandler = (*consumerGroup)(nil)

type (
	ConsumeBatchHandler func(msgs []*sarama.ConsumerMessage) (markFlags []bool, errs []error)
	consumerGroup       struct {
		Ready        chan bool
		groupId      string
		maxInterval  time.Duration
		batchSize    int
		batchHandler ConsumeBatchHandler

		client sarama.ConsumerGroup
	}
)

func newConsumerGroup(addrs []string, version sarama.KafkaVersion, groupId string, maxInterval time.Duration, batchSize int) (*consumerGroup, error) {
	if groupId == "" {
		return nil, errors.New("groupId must not be empty")
	}
	if maxInterval == 0 || batchSize == 0 {
		return nil, errors.New("maxInterval or batchSize must not be zero")
	}
	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.IsolationLevel = sarama.ReadCommitted
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.Retry.Max = 3
	config.Consumer.Retry.Backoff = time.Second * 1
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}

	client, err := sarama.NewConsumerGroup(addrs, groupId, config)
	if err != nil {
		return nil, err
	}
	consumer := &consumerGroup{
		Ready:       make(chan bool),
		groupId:     groupId,
		maxInterval: maxInterval,
		batchSize:   batchSize,
		client:      client,
	}
	return consumer, nil
}

func (c *consumerGroup) handleMessages(session sarama.ConsumerGroupSession, messages []*sarama.ConsumerMessage) {
	markFlags, errs := c.batchHandler(messages)

	var (
		consumed   int
		notNilErrs []error
	)
	for i, msg := range messages {
		if markFlags[i] {
			session.MarkMessage(msg, "")
			consumed++
		}
		if errs == nil {
			continue
		}
		if err := errs[i]; err != nil {
			notNilErrs = append(notNilErrs, err)
		}
	}
	session.Commit()

	if consumed > 0 {
		logx.Debug(c.groupId, "consumed message: %d", consumed)
	}
	if notNilErrs != nil {
		logx.Error(c.groupId, "failed handleBatch error count: %d, msg: %s", len(notNilErrs),
			errors.Join(notNilErrs...).Error())
	}
}

func (c *consumerGroup) Setup(sarama.ConsumerGroupSession) error {
	close(c.Ready)
	return nil
}

func (c *consumerGroup) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	messages := make([]*sarama.ConsumerMessage, 0, c.batchSize)

	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				err := errors.New("message channel was closed")
				logx.Error(c.groupId, err.Error())
				return err
			}
			messages = append(messages, msg)

			if len(messages) >= c.batchSize {
				c.handleMessages(session, messages)
				messages = messages[:0]
			}
		case <-time.After(c.maxInterval):
			if len(messages) > 0 {
				c.handleMessages(session, messages)
				messages = messages[:0]
			}
		case <-session.Context().Done():
			return nil
		}
	}
}
