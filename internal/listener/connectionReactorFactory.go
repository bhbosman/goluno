package listener

import (
	"context"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/log"
	"github.com/cskr/pubsub"
)

type ConnectionReactorFactory struct {
	name          string
	PubSub        *pubsub.PubSub
	SerializeData SerializeData
}

func (self *ConnectionReactorFactory) Create(
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	logger *log.SubSystemLogger,
	userContext interface{}) commsImpl.IConnectionReactor {
	result := NewConnectionReactor(
		logger,
		name,
		cancelCtx,
		cancelFunc,
		userContext,
		self.PubSub,
		self.SerializeData)
	return result
}

func (self *ConnectionReactorFactory) Name() string {
	return self.name
}

func NewConnectionReactorFactory(
	name string,
	pubSub *pubsub.PubSub,
	SerializeData SerializeData) *ConnectionReactorFactory {
	return &ConnectionReactorFactory{
		name:          name,
		PubSub:        pubSub,
		SerializeData: SerializeData,
	}
}
