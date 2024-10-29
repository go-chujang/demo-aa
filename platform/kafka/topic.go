package kafka

import (
	"fmt"
)

var (
	TopicMembership = Topic{topicName: "membership"}
	TopicLiquidity  = Topic{topicName: "liquidity"}
	TopicOperation  = Topic{topicName: "operation"}
)

type Topic struct {
	topicName string
}

func (t Topic) String() string      { return t.topicName }
func (t Topic) Equal(s string) bool { return t.topicName == s }

func TopicName(topic string) (Topic, error) {
	switch topic {
	case TopicMembership.String():
		return TopicMembership, nil
	case TopicLiquidity.String():
		return TopicLiquidity, nil
	case TopicOperation.String():
		return TopicOperation, nil
	default:
		return Topic{}, fmt.Errorf("not exist topic: %s", topic)
	}
}
