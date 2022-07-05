package listener

import (
	"context"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type Factory struct {
	PubSub          *pubsub.PubSub
	SerializeData   SerializeData
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl
}

func (self *Factory) Create(
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
		userContext,
		self.PubSub,
		self.SerializeData,
		self.ConsumerCounter)
	return result, nil
}

func NewConnectionReactorFactory(
	pubSub *pubsub.PubSub,
	SerializeData SerializeData,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl,
) (*Factory, error) {
	fac := &Factory{
		PubSub:          pubSub,
		SerializeData:   SerializeData,
		ConsumerCounter: ConsumerCounter,
	}
	return fac, nil
}
