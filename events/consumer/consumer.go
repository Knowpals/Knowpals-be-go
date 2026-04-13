package consumer

import (
	"context"

	"github.com/IBM/sarama"
)

type Consumer interface {
	Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error
}

type saramaConsumer struct {
	consumerGroup sarama.ConsumerGroup
}

func NewSaramaConsumer(client sarama.Client, groupID string) Consumer {
	cg, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		panic("创建消费者失败")
	}
	return &saramaConsumer{consumerGroup: cg}
}

func (c *saramaConsumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	return c.consumerGroup.Consume(ctx, topics, handler)
}
