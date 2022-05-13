package pubsub

import (
	"github.com/alash3al/go-pubsub"
)

var DefaultBroker = Build()

func Build() *PubSub {
	return &PubSub{
		pubsub.NewBroker(),
		make(map[string]*pubsub.Subscriber, 10),
	}
}

type PubSub struct {
	*pubsub.Broker
	subscribers map[string]*pubsub.Subscriber
}

func(p *PubSub)GetSubscriber(brokerName, topic string) *pubsub.Subscriber {
	var subscriber *pubsub.Subscriber
	var ok bool
	if subscriber, ok = p.subscribers[brokerName]; !ok{
		// ok 直接拿
		p.subscribers[brokerName], _ = p.Broker.Attach()
		subscriber = p.subscribers[brokerName]
	}
	p.Broker.Subscribe(subscriber,topic)
	return subscriber
}