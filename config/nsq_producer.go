package config

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
)

type NSQProducer struct {
	producer *nsq.Producer
}

func newNSQProducer(address string) (*NSQProducer, error) {
	producerConfig := nsq.NewConfig()
	producer, err := nsq.NewProducer(address, producerConfig)
	if err != nil {
		return nil, err
	}

	return &NSQProducer{
		producer: producer,
	}, nil
}

func (np *NSQProducer) Publish(topic string, message interface{}) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return np.producer.Publish(topic, jsonMessage)
}
