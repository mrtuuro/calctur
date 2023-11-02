package main

import (
	"encoding/json"
	"time"

	"calctur/backend/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DataProducer interface {
	ProduceData(data types.Coordinate) error
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					//fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					//fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &KafkaProducer{
		producer: p,
		topic:    topic,
	}, nil
}

func (kf *KafkaProducer) ProduceData(data types.Coordinate) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return kf.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kf.topic,
			Partition: kafka.PartitionAny,
		},
		Value:     b,
		Timestamp: time.Now().UTC(),
	}, nil)
}
