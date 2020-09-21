package lunoWS

import (
	"context"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"github.com/cskr/pubsub"
)

type ConnectionReactorFactory struct {
	name         string
	APIKeyID     string
	APIKeySecret string
	PubSub       *pubsub.PubSub
}

func (self *ConnectionReactorFactory) Create(
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	logger *gologging.SubSystemLogger,
	userContext interface{}) intf.IConnectionReactor {
	result := NewConnectionReactor(
		logger,
		name,
		cancelCtx,
		cancelFunc,
		self.APIKeyID,
		self.APIKeySecret,
		self.PubSub,
		userContext)
	return result
}

func (self *ConnectionReactorFactory) Name() string {
	return self.name
}

func NewConnectionReactorFactory(
	name string,
	APIKeyID string,
	APIKeySecret string,
	pubSub *pubsub.PubSub) *ConnectionReactorFactory {
	return &ConnectionReactorFactory{
		name:         name,
		APIKeyID:     APIKeyID,
		APIKeySecret: APIKeySecret,
		PubSub:       pubSub,
	}
}
