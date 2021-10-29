package lunoWS

import (
	"context"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type ConnectionReactorFactory struct {
	APIKeyID     string
	APIKeySecret string
	PubSub       *pubsub.PubSub
	name         string
}

func (self *ConnectionReactorFactory) Name() string {
	return self.name
}

func (self *ConnectionReactorFactory) Values(inputValues map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	return result, nil
}

func (self *ConnectionReactorFactory) Create(
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc common.ConnectionCancelFunc,
	logger *zap.Logger,
	userContext interface{}) intf.IConnectionReactor {
	result := NewConnectionReactor(
		logger,
		cancelCtx,
		cancelFunc,
		connectionCancelFunc,
		self.APIKeyID,
		self.APIKeySecret,
		self.PubSub,
		userContext)
	return result
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
