package kafka

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

var (
	producerPool *producerProvider
	producerOnce sync.Once
	consumer     *consumerGroup
	consumerOnce sync.Once
)

func producerMustNotNil() {
	if producerPool == nil {
		panic("producer must not nil, UseProducer() first")
	}
}

func consumerMustNotNil() {
	if consumer == nil {
		panic("consumer must not nil, UseConsumer() first")
	}
}

func UseProducer(addrs []string, version sarama.KafkaVersion) (err error) {
	producerOnce.Do(func() { producerPool, err = newProducerProvider(addrs, version) })
	return err
}

func UseConsumer(addrs []string, version sarama.KafkaVersion, groupId string, maxInterval time.Duration, batchSize int) (err error) {
	consumerOnce.Do(func() { consumer, err = newConsumerGroup(addrs, version, groupId, maxInterval, batchSize) })
	return err
}

func Unuse() (err error) {
	if producerPool != nil {
		producerPool.clear()
	}
	if consumer != nil {
		err = consumer.client.Close()
	}
	return err
}

func Version(v string) (sarama.KafkaVersion, error) { return sarama.ParseKafkaVersion(v) }

func EnvAddrsVersion() (addrs []string, version sarama.KafkaVersion) {
	addrs = strings.Split(os.Getenv("KAFKA_ADDRS"), ",")
	version, _ = Version(os.Getenv("KAFKA_VERSION"))
	return addrs, version
}
