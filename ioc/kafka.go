package ioc

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/Knowpals/Knowpals-be-go/config"
)

func InitKafka(cfg *config.Config) sarama.Client {
	saramaCfg := sarama.NewConfig()

	saramaCfg.Net.SASL.Enable = false
	saramaCfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	client, err := sarama.NewClient(cfg.Kafka.Addrs, saramaCfg)
	if err != nil {
		log.Fatal("初始化 kafka 失败", err)
	}
	return client
}

func InitKafkaConsumerGroupID(cfg *config.Config) string {
	return cfg.Kafka.ConsumerGroup
}
