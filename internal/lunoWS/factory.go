package lunoWS

import (
	"context"
	"github.com/bhbosman/gocommon/model"
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

func (self *ConnectionReactorFactory) Values(_ map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	return result, nil
}

func (self *ConnectionReactorFactory) Create(
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	userContext interface{}) (intf.IConnectionReactor, error) {
	result := NewConnectionReactor(
		logger,
		cancelCtx,
		cancelFunc,
		connectionCancelFunc,
		self.APIKeyID,
		self.APIKeySecret,
		self.PubSub,
		userContext)
	return result, nil
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
