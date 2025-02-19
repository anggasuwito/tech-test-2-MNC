package consumer

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type subscriber struct {
	Topic   string
	Channel string
	Handler func(message *nsq.Message) error
	Option  *SubscriberOption
}

type NsqConsumer struct {
	Connection  string
	MaxInFlight int
	RequeueTime int
	MaxAttempts uint16
}

type SubscriberOption struct {
	ConfigKey   string
	Nconsumer   int
	MaxInFlight int
	RequeueTime int
	MaxAttempts uint16
}

type NSQ struct {
	Connection    []string
	Consumers     map[string]*subscriber
	ListenErrCh   chan error
	stopCh        chan struct{}
	DefaultOption NsqConsumer
}

func New(connection []string) *NSQ {
	nsqClient := &NSQ{
		Consumers:   make(map[string]*subscriber),
		ListenErrCh: make(chan error),
		stopCh:      make(chan struct{}),
		Connection:  connection,
	}

	return nsqClient
}

func (n *NSQ) Run() {
	for _, consumer := range n.Consumers {
		go func(consumer *subscriber) {
			maxAttempts := consumer.Option.MaxAttempts
			maxInFlight := consumer.Option.MaxInFlight
			requeueTime := consumer.Option.RequeueTime
			nConsumer := consumer.Option.Nconsumer

			nsqConfig := nsq.NewConfig()
			nsqConfig.MaxInFlight = maxInFlight
			nsqConfig.MaxAttempts = maxAttempts

			q, _ := nsq.NewConsumer(consumer.Topic, consumer.Channel, nsqConfig)
			log.Printf("[!!!] RUNNING CONSUMER (t: %s), (n: %d), (maxInFlight: %d), (maxRetryAttempts: %d), (requeueTime: %v)\n", consumer.Topic, nConsumer, maxInFlight, maxAttempts, requeueTime)
			q.AddConcurrentHandlers(nsq.HandlerFunc(consumer.Handler), nConsumer)

			err := q.ConnectToNSQLookupds(n.Connection)
			if err != nil {
				n.ListenErrCh <- err
			}

			<-n.stopCh
			q.Stop()

		}(consumer)
	}

	// Wait for interrupt signal to stop the consumers gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Stopping consumers gracefully...")

	// Send stop signal to all consumers
	close(n.stopCh)

	// Wait for all consumers to stop
	for _, consumer := range n.Consumers {
		log.Printf(fmt.Sprintf("CONSUMER STOPPED Topic:%s Channel:%s", consumer.Topic, consumer.Channel))
		<-n.stopCh
	}
}

func (n *NSQ) Register(topic, channel string, handler func(message *nsq.Message) error, option *SubscriberOption) {
	log.Printf("Registering NSQ Consumer : %s/%s", topic, channel)

	configKey := fmt.Sprintf("%s/%s", topic, channel)
	n.Consumers[configKey] = &subscriber{
		Topic:   topic,
		Channel: channel,
		Handler: handler,
		Option:  option,
	}
}

func (n *NSQ) NewSubscriberOption(configKey string, opt *struct {
	Nconsumer   int
	MaxInFlight int
	RequeueTime int
	MaxAttempts uint16
}) *SubscriberOption {
	return &SubscriberOption{
		ConfigKey:   configKey,
		Nconsumer:   opt.Nconsumer,
		MaxInFlight: opt.MaxInFlight,
		RequeueTime: opt.RequeueTime,
		MaxAttempts: opt.MaxAttempts,
	}
}

func (n *NSQ) GetSubscriberOptions(configKey string) (option *SubscriberOption) {
	option = &SubscriberOption{
		ConfigKey:   configKey,
		MaxInFlight: n.DefaultOption.MaxInFlight,
		RequeueTime: n.DefaultOption.RequeueTime,
		MaxAttempts: n.DefaultOption.MaxAttempts,
	}
	if consumer, ok := n.Consumers[configKey]; ok {
		option = consumer.Option
	}
	return
}

func (n *NSQ) GetConsumerRequeueTime(configKey string) (requeueTime time.Duration) {
	return time.Duration(n.GetSubscriberOptions(configKey).RequeueTime) * time.Second
}

func (n *NSQ) GetConsumerMaxAttempts(configKey string) (maxAttempts uint16) {
	return n.GetSubscriberOptions(configKey).MaxAttempts
}

func (n *NSQ) GetConsumerMaxInFlight(configKey string) (maxInFlight int) {
	return n.GetSubscriberOptions(configKey).MaxInFlight
}

func (n *NSQ) GetNumberOfConsumer(configKey string) (nConsumer int) {
	return n.GetSubscriberOptions(configKey).Nconsumer
}
