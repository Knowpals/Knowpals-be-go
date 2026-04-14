package ioc

import (
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/Knowpals/Knowpals-be-go/config"
)

func InitKafka(cfg *config.Config) sarama.Client {
	saramaCfg := sarama.NewConfig()

	saramaCfg.Version = sarama.V2_1_0_0

	saramaCfg.Net.SASL.Enable = false
	saramaCfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner

	// consumer group（关键）
	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	saramaCfg.Consumer.Group.Session.Timeout = 60 * time.Second
	saramaCfg.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	saramaCfg.Consumer.MaxProcessingTime = 5 * time.Minute
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaCfg.Consumer.Return.Errors = true
	client, err := sarama.NewClient(cfg.Kafka.Addrs, saramaCfg)
	if err != nil {
		log.Fatal("初始化 kafka 失败", err)
	}
	return client
}

func InitKafkaConsumerGroupID(cfg *config.Config) string {
	return cfg.Kafka.ConsumerGroup
}
