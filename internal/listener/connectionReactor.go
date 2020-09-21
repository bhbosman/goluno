package listener

import (
	"context"
	"github.com/bhbosman/goLuno/internal/ConsumerCounter"
	"github.com/bhbosman/goLuno/internal/common"
	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/proto"
	"net"
	"net/url"
)

type SerializeData func(m proto.Message) (goprotoextra.IReadWriterSize, error)
type ConnectionReactor struct {
	impl.BaseConnectionReactor
	messageRouter *messageRouter.MessageRouter
	PubSub        *pubsub.PubSub
	Pairs         []*common.PairInformation
	SerializeData SerializeData
	ConsumerCounter *ConsumerCounter.ConsumerCounter
}

func NewConnectionReactor(
	logger *gologging.SubSystemLogger,
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	userContext interface{},
	PubSub *pubsub.PubSub,
	SerializeData SerializeData,
	ConsumerCounter *ConsumerCounter.ConsumerCounter) *ConnectionReactor {
	Pairs, _ := userContext.([]*common.PairInformation)
	result := &ConnectionReactor{
		BaseConnectionReactor: impl.NewBaseConnectionReactor(logger, name, cancelCtx, cancelFunc, userContext),
		messageRouter:         messageRouter.NewMessageRouter(),
		PubSub:                PubSub,
		Pairs:                 Pairs,
		SerializeData:         SerializeData,
		ConsumerCounter:       ConsumerCounter,
	}
	return result
}

func (self *ConnectionReactor) Init(
	conn net.Conn,
	url *url.URL,
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	toConnectionFunc goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc) (rxgo.NextExternalFunc, error) {
	_, err := self.BaseConnectionReactor.Init(conn, url, connectionId, connectionManager, toConnectionFunc, toConnectionReactor)
	if err != nil {
		return nil, err
	}
	self.messageRouter.Add(self.HandleTop5)

	var republishTopics []string
	var publishTopics []string
	for _, pair := range self.Pairs {
		republishTopics = append(republishTopics, common.RepublishName(pair.Pair))
		publishTopics = append(publishTopics, common.PublishName(pair.Pair))
	}
	ch := self.PubSub.Sub(publishTopics...)
	go func(ch chan interface{}, topics ...string) {
		defer self.PubSub.Unsub(ch, topics...)
		for {
			select {
			case <-self.CancelCtx.Done():
				return
			case v := <-ch:

				self.ToReactor(false, v)
			}
		}
	}(ch, publishTopics...)

	self.PubSub.Pub(&struct{}{}, republishTopics...)

	return self.doNext, nil
}

func (self *ConnectionReactor) doNext(external bool, i interface{}) {
	_, _ = self.messageRouter.Route(i)
}

func (self *ConnectionReactor) HandleTop5(top5 *marketDataStream.PublishTop5) error {
	marshal, err := self.SerializeData(top5)
	if err != nil {
		return err
	}
	_ = self.ToConnection(marshal)
	return nil
}


func (self *ConnectionReactor) Open() error {
	self.ConsumerCounter.AddConsumer()
	return self.BaseConnectionReactor.Open()
}


func (self *ConnectionReactor) Close() error {
	self.ConsumerCounter.RemoveConsumer()
	return self.BaseConnectionReactor.Close()
}
