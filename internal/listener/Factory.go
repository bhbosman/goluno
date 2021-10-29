package listener

import (
	"context"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type Factory struct {
	name            string
	PubSub          *pubsub.PubSub
	SerializeData   SerializeData
	ConsumerCounter *netDial.CanDialDefaultImpl
}

func (self *Factory) Values(inputValues map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	return result, nil
}

func (self *Factory) Create(
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
		userContext,
		self.PubSub,
		self.SerializeData,
		self.ConsumerCounter)
	return result
}

func (self *Factory) Name() string {
	return self.name
}

func NewConnectionReactorFactory(
	name string,
	pubSub *pubsub.PubSub,
	SerializeData SerializeData,
	ConsumerCounter *netDial.CanDialDefaultImpl) intf.IConnectionReactorFactory {
	return &Factory{
		name:            name,
		PubSub:          pubSub,
		SerializeData:   SerializeData,
		ConsumerCounter: ConsumerCounter,
	}
}
