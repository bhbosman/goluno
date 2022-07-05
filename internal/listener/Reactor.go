package listener

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goLuno/internal/common"
	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messageRouter"
	common3 "github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"strings"
)

type SerializeData func(m proto.Message) (goprotoextra.IReadWriterSize, error)
type Reactor struct {
	common2.BaseConnectionReactor
	messageRouter   *messageRouter.MessageRouter
	PubSub          *pubsub.PubSub
	Pairs           []*common.PairInformation
	SerializeData   SerializeData
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl
}

func (self *Reactor) Init(
	toConnectionFunc goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc,
	toConnectionFuncReplacement rxgo.NextFunc,
	toConnectionReactorReplacement rxgo.NextFunc,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	_, _, _, err := self.BaseConnectionReactor.Init(
		toConnectionFunc,
		toConnectionReactor,
		toConnectionFuncReplacement,
		toConnectionReactorReplacement,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	_ = self.messageRouter.Add(self.HandleTop5)

	var republishTopics []string
	var publishTopics []string
	for _, pair := range self.Pairs {
		republishTopics = append(republishTopics, common.RepublishName(pair.Pair))
		publishTopics = append(publishTopics, common.PublishName(pair.Pair))
	}
	if publishTopics != nil {
		ch := self.PubSub.Sub(publishTopics...)
		go func(ch chan interface{}, topics ...string) {

			<-self.CancelCtx.Done()
			self.PubSub.Unsub(ch, topics...)
		}(ch, publishTopics...)

		go func(ch chan interface{}, topics ...string) {
			for v := range ch {
				if self.CancelCtx.Err() == nil {
					_ = self.ToReactor(false, v)
				}
			}
		}(ch, publishTopics...)

		self.PubSub.Pub(&struct{}{}, republishTopics...)
	}

	return func(i interface{}) {
			self.doNext(false, i)
		},
		func(err error) {
			self.doNext(false, err)
		},
		func() {

		}, nil
}

func (self *Reactor) doNext(_ bool, i interface{}) {
	_, _ = self.messageRouter.Route(i)
}

func (self *Reactor) HandleTop5(top5 *marketDataStream.PublishTop5) error {
	if self.CancelCtx.Err() != nil {
		return self.CancelCtx.Err()
	}
	top5.Source = "LunoWS"
	s := strings.Replace(fmt.Sprintf("%v.%v", top5.Source, top5.Instrument), "/", ".", -1)
	top5.UniqueName = s
	if self.SerializeData != nil {
		marshal, err := self.SerializeData(top5)
		if err != nil {
			return err
		}
		_ = self.ToConnection(marshal)
	}
	return nil
}

func (self *Reactor) Open() error {
	self.ConsumerCounter.AddConsumer()
	return self.BaseConnectionReactor.Open()
}

func (self *Reactor) Close() error {
	self.ConsumerCounter.RemoveConsumer()
	return self.BaseConnectionReactor.Close()
}

func NewConnectionReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc common3.ConnectionCancelFunc,
	userContext interface{},
	PubSub *pubsub.PubSub,
	SerializeData SerializeData,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl) *Reactor {
	Pairs, _ := userContext.([]*common.PairInformation)
	result := &Reactor{
		BaseConnectionReactor: common2.NewBaseConnectionReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			userContext),
		messageRouter:   messageRouter.NewMessageRouter(),
		PubSub:          PubSub,
		Pairs:           Pairs,
		SerializeData:   SerializeData,
		ConsumerCounter: ConsumerCounter,
	}
	return result
}
