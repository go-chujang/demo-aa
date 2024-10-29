package main

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/go-chujang/demo-aa/platform/kafka"
)

func setup_kafka() {
	deleteTopics()
	time.Sleep(time.Second * 1)

	for topicName, numPartitions := range map[kafka.Topic]int32{
		kafka.TopicMembership: 1,
		kafka.TopicLiquidity:  1,
		kafka.TopicOperation:  2,
	} {
		if err = createTopicByEnv(topicName, numPartitions); err != nil {
			panic(err)
		}
	}
}

func deleteTopics() {
	addrs, version := kafka.EnvAddrsVersion()
	cfg := sarama.NewConfig()
	cfg.Version = version

	admin, err := sarama.NewClusterAdmin(addrs, cfg)
	if err != nil {
		panic(err)
	}
	defer admin.Close()
	admin.DeleteTopic(kafka.TopicMembership.String())
	admin.DeleteTopic(kafka.TopicLiquidity.String())
	admin.DeleteTopic(kafka.TopicOperation.String())
}

func createTopicByEnv(topic kafka.Topic, numPartitions int32) error {
	addrs, version := kafka.EnvAddrsVersion()
	cfg := sarama.NewConfig()
	cfg.Version = version

	admin, err := sarama.NewClusterAdmin(addrs, cfg)
	if err != nil {
		return err
	}
	defer admin.Close()
	return admin.CreateTopic(topic.String(), &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: 1,
	}, false)
}
