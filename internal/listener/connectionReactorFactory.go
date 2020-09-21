package listener

import (
	"context"
	"github.com/bhbosman/goLuno/internal/ConsumerCounter"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gologging"
	"github.com/cskr/pubsub"
)

type ConnectionReactorFactory struct {
	name          string
	PubSub        *pubsub.PubSub
	SerializeData SerializeData
	ConsumerCounter *ConsumerCounter.ConsumerCounter
}

func (self *ConnectionReactorFactory) Create(
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	logger *gologging.SubSystemLogger,
	userContext interface{}) commsImpl.IConnectionReactor {
	result := NewConnectionReactor(
		logger,
		name,
		cancelCtx,
		cancelFunc,
		userContext,
		self.PubSub,
		self.SerializeData,
		self.ConsumerCounter)
	return result
}

func (self *ConnectionReactorFactory) Name() string {
	return self.name
}

func NewConnectionReactorFactory(
	name string,
	pubSub *pubsub.PubSub,
	SerializeData SerializeData,
	ConsumerCounter *ConsumerCounter.ConsumerCounter) *ConnectionReactorFactory {
	return &ConnectionReactorFactory{
		name:            name,
		PubSub:          pubSub,
		SerializeData:   SerializeData,
		ConsumerCounter: ConsumerCounter,
	}
}
