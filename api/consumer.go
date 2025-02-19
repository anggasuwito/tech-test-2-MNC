package api

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	nsqd "tech-test-2-MNC/pkg/nsq"

	"tech-test-2-MNC/config"
	"tech-test-2-MNC/internal/constant"
	"tech-test-2-MNC/internal/handler/consumer"
	"tech-test-2-MNC/internal/repository"
	"tech-test-2-MNC/internal/usecase"
)

type NSQ struct {
	client             *nsqd.NSQ
	transactionHandler consumer.TransactionHandler
}

func NewConsumer() *NSQ {
	cfg := config.GetConfig()
	nsqClient := nsqd.New(cfg.NSQConsumerHost)

	txWrapper := repository.NewTransactionWrapper(cfg.DBMaster)
	accountRepo := repository.NewUserAccountRepo(cfg.DBMaster)
	transactionRepo := repository.NewTransactionRepo(cfg.DBMaster)

	tuc := usecase.NewTransactionUC(txWrapper, accountRepo, transactionRepo)
	shHandler := consumer.NewTransactionHandler(nsqClient, tuc)

	return &NSQ{
		client:             nsqClient,
		transactionHandler: shHandler,
	}
}

func (n *NSQ) RegisterAll() *NSQ {
	cfg := config.GetConfig()
	consumerList := map[string]func(message *nsq.Message) error{
		constant.ConsumerUpdateTransactionStatus: n.transactionHandler.UpdateTransactionStatus,
	}

	for _, subscriber := range cfg.NSQConsumers {
		configKey := fmt.Sprintf("%s/%s", subscriber.Topic, subscriber.Channel)
		opt := n.client.NewSubscriberOption(configKey, &struct {
			Nconsumer   int
			MaxInFlight int
			RequeueTime int
			MaxAttempts uint16
		}{
			Nconsumer:   cfg.NSQNConsumer,
			MaxInFlight: cfg.NSQMaxInFlight,
			RequeueTime: cfg.NSQRequeueTime,
			MaxAttempts: cfg.NSQMaxAttempts,
		})

		if _, exist := consumerList[configKey]; !exist {
			log.Println(fmt.Sprintf("[Consumer][RegisterAll] Topic : %s NOT FOUND", configKey))
			continue
		}

		n.client.Register(subscriber.Topic, subscriber.Channel, consumerList[configKey], opt)
	}
	return n
}

func (n *NSQ) Run() {
	n.client.Run()
}
