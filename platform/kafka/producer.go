package kafka

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-chujang/demo-aa/common/utils/id"
)

type producerProvider struct {
	mu               sync.Mutex
	configProvider   func() *sarama.Config
	producerProvider func() (sarama.AsyncProducer, error)
	producers        []sarama.AsyncProducer
}

func newProducerProvider(addrs []string, version sarama.KafkaVersion) (*producerProvider, error) {
	provider := &producerProvider{}
	provider.configProvider = func() *sarama.Config {
		config := sarama.NewConfig()
		config.Version = version
		config.Net.MaxOpenRequests = 1
		config.Producer.Idempotent = true
		config.Producer.Return.Errors = false
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
		config.Producer.Transaction.ID = id.Uuid()
		config.Producer.Transaction.Retry.Backoff = 10
		config.Producer.Transaction.Retry.Max = 3
		config.Producer.Transaction.Timeout = time.Second * 3
		return config
	}
	provider.producerProvider = func() (sarama.AsyncProducer, error) {
		id := id.Seq()
		config := provider.configProvider()
		config.Producer.Transaction.ID = fmt.Sprintf("%s-%d", config.Producer.Transaction.ID, id)
		producer, err := sarama.NewAsyncProducer(addrs, config)
		if err != nil {
			return nil, err
		}
		return producer, nil
	}
	test, err := provider.borrow()
	if err != nil {
		return nil, err
	}
	provider.release(test)
	return provider, nil
}

func (p *producerProvider) borrow() (sarama.AsyncProducer, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.producers) == 0 {
		var errs []error
		for i := 0; i < 3; i++ {
			producer, err := p.producerProvider()
			if err == nil {
				// success
				return producer, nil
			}
			errs = append(errs, err)
		}
		return nil, errors.Join(errs...)
	}

	index := len(p.producers) - 1
	producer := p.producers[index]
	p.producers = p.producers[:index]
	return producer, nil
}

func (p *producerProvider) release(producer sarama.AsyncProducer) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if producer.TxnStatus()&sarama.ProducerTxnFlagInError != 0 {
		_ = producer.Close()
		return
	}
	p.producers = append(p.producers, producer)
}

func (p *producerProvider) clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, producer := range p.producers {
		producer.Close()
	}
	p.producers = p.producers[:0]
}
