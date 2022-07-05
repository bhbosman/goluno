package lunoWS

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/fullMarketData"
	"github.com/bhbosman/goLuno/internal/lunoWS/internal"
	lunaRawDataFeed "github.com/bhbosman/goMessages/luno/stream"
	marketDataStream "github.com/bhbosman/goMessages/marketData/stream"
	"github.com/bhbosman/gocommon/messages"
	common3 "github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gomessageblock"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"

	"github.com/bhbosman/gocommon/messageRouter"

	"github.com/bhbosman/goCommsStacks/webSocketMessages/wsmsg"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type Reactor struct {
	common2.BaseConnectionReactor
	messageRouter        *messageRouter.MessageRouter
	APIKeyID             string
	APIKeySecret         string
	FullMarketOrderBook  *fullMarketData.FullMarketOrderBook
	PubSub               *pubsub.PubSub
	LunaPairInformation  *common.PairInformation
	republishChannelName string
	publishChannelName   string
	UpdateCount          int64
	sequence             int64
}

func (self *Reactor) SendMessage(message proto.Message) error {
	rws := gomessageblock.NewReaderWriter()
	m := jsonpb.Marshaler{}
	err := m.Marshal(rws, message)
	if err != nil {
		return err
	}
	flatten, err := rws.Flatten()
	if err != nil {
		return err
	}
	WebSocketMessage := wsmsg.WebSocketMessage{
		OpCode:  wsmsg.WebSocketMessage_OpText,
		Message: flatten,
	}
	readWriterSize, err := stream.Marshall(&WebSocketMessage)
	if err != nil {
		return err
	}

	return self.ToConnection(readWriterSize)
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
	_ = self.messageRouter.Add(self.HandleWebSocketMessageWrapper)
	_ = self.messageRouter.Add(self.HandleReaderWriter)
	_ = self.messageRouter.Add(self.HandleLunaData)
	_ = self.messageRouter.Add(self.HandleEmptyQueue)
	_ = self.messageRouter.Add(self.HandlePublishMessage)

	self.republishChannelName = common.RepublishName(self.LunaPairInformation.Pair)
	self.publishChannelName = common.PublishName(self.LunaPairInformation.Pair)

	republishChannel := self.PubSub.Sub(self.republishChannelName)
	go func(ch chan interface{}, topics ...string) {
		<-self.CancelCtx.Done()
		self.PubSub.Unsub(ch, topics...)
	}(republishChannel, self.republishChannelName)

	go func(ch chan interface{}, topics ...string) {
		for range ch {
			if self.CancelCtx.Err() == nil {
				_ = self.ToReactor(false, &internal.PublishMessage{})
			}
		}
	}(republishChannel, self.republishChannelName)

	return func(i interface{}) {
			self.doNext(false, i)
		},
		func(err error) {
			self.doNext(false, err)
		},
		func() {
			d := 23
			d++
		}, nil
}

func (self *Reactor) Close() error {
	top5 := &marketDataStream.PublishTop5{
		Instrument: self.LunaPairInformation.Pair,
	}
	self.PubSub.Pub(top5, self.publishChannelName)

	err := self.BaseConnectionReactor.Close()
	return err
}

func (self *Reactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	return err
}

func (self *Reactor) doNext(_ bool, i interface{}) {
	_, err := self.messageRouter.Route(i)
	if err != nil {
		return
	}
}

func (self *Reactor) HandleReaderWriter(msg *gomessageblock.ReaderWriter) {
	marshal, err := stream.UnMarshal(
		msg,
		func(i interface{}) {
			self.ToReactor(false, i)
		},
		func(i interface{}) {
			self.ToConnection(msg)
		},
	)
	if err != nil {
		//return err
	}
	_, err = self.messageRouter.Route(marshal)
	if err != nil {
		//return err
	}
}

func (self *Reactor) HandlePublishMessage(_ *internal.PublishMessage) error {
	return self.publishData(true)
}

func (self *Reactor) HandleEmptyQueue(_ *messages.EmptyQueue) error {
	return self.publishData(false)
}
func (self *Reactor) publishData(forcePublish bool) error {
	thereWasAChange := forcePublish
	var bids []*marketDataStream.Point
	if highBidNode := self.FullMarketOrderBook.OrderSide[fullMarketData.BuySide].Right(); highBidNode != nil {
		count := 0
		for node := highBidNode; node != nil; node = node.Prev() {
			bidPrice := node.Key.(float64)
			if pp, ok := node.Value.(*fullMarketData.PricePoint); ok {
				thereWasAChange = thereWasAChange || pp.Touched
				pp.ClearTouched()
				bids = append(bids, &marketDataStream.Point{
					Price:  bidPrice,
					Volume: pp.GetVolume(),
				})
			}
			count++
		}
	}
	var asks []*marketDataStream.Point
	if lowAskNode := self.FullMarketOrderBook.OrderSide[fullMarketData.AskSide].Left(); lowAskNode != nil {
		count := 0
		for node := lowAskNode; node != nil; node = node.Next() {
			askPrice := node.Key.(float64)
			if pp, ok := node.Value.(*fullMarketData.PricePoint); ok {
				thereWasAChange = thereWasAChange || pp.Touched
				pp.ClearTouched()
				asks = append(asks, &marketDataStream.Point{
					Price:  askPrice,
					Volume: pp.GetVolume(),
				})
			}
			count++
		}
	}
	spread := 0.0
	if len(asks) > 0 && len(bids) > 0 {
		spread = asks[0].Price - bids[0].Price
	}
	if thereWasAChange {
		if !forcePublish {
			self.UpdateCount++
		}
		top5 := &marketDataStream.PublishTop5{
			Instrument:         self.LunaPairInformation.Pair,
			Spread:             spread,
			SourceTimeStamp:    self.FullMarketOrderBook.SourceTimestamp,
			SourceMessageCount: self.FullMarketOrderBook.SourceMessageCount,
			UpdateCount:        self.UpdateCount,
			Bid:                bids,
			Ask:                asks,
		}
		self.PubSub.Pub(top5, self.publishChannelName)
	}

	return nil
}
func (self *Reactor) updateSequence(newSeq int64) error {
	if self.sequence == 0 {
		self.sequence = newSeq
		return nil
	}
	if self.sequence+1 == newSeq {
		self.sequence = newSeq
		return nil
	}
	return fmt.Errorf("invalid sequence. Expected: %v, received: %v", self.sequence+1, newSeq)
}

func (self *Reactor) HandleLunaData(msg *lunaRawDataFeed.LunoStreamData) error {
	err := self.updateSequence(msg.Sequence)
	if err != nil {
		self.CancelFunc()
		return err
	}
	self.FullMarketOrderBook.SetTimeStamp(msg.Timestamp)
	self.FullMarketOrderBook.UpdateMessageReceivedCount()
	switch {
	case msg.Bids != nil || msg.Asks != nil:
		self.FullMarketOrderBook.Clear()
		for _, order := range msg.Bids {
			self.FullMarketOrderBook.AddOrder(fullMarketData.BuySide, order)
		}
		for _, order := range msg.Asks {
			self.FullMarketOrderBook.AddOrder(fullMarketData.AskSide, order)
		}
	case msg.TradeUpdates != nil:
		for _, order := range msg.TradeUpdates {
			self.FullMarketOrderBook.TradeUpdate(order)
		}
	case msg.DeleteUpdate != nil:
		self.FullMarketOrderBook.DeleteUpdate(msg.DeleteUpdate)

	case msg.CreateUpdate != nil:
		self.FullMarketOrderBook.CreateUpdate(msg.CreateUpdate)
	}
	return nil
}

func (self *Reactor) HandleWebSocketMessageWrapper(msg *wsmsg.WebSocketMessageWrapper) error {
	switch msg.Data.OpCode {
	case wsmsg.WebSocketMessage_OpText:
		Unmarshaler := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
			AnyResolver:        nil,
		}
		lunaData := &lunaRawDataFeed.LunoStreamData{}
		err := Unmarshaler.Unmarshal(bytes.NewBuffer(msg.Data.Message), lunaData)
		if err != nil {
			return err
		}
		_, _ = self.messageRouter.Route(lunaData)
		return nil
	case wsmsg.WebSocketMessage_OpStartLoop:
		msg := &lunaRawDataFeed.Credentials{
			ApiKeyId:     self.APIKeyID,
			ApiKeySecret: self.APIKeySecret,
		}
		_ = self.SendMessage(msg)
		return nil
	default:
		return nil
	}
}

func NewConnectionReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc common3.ConnectionCancelFunc,
	APIKeyID string,
	APIKeySecret string,
	PubSub *pubsub.PubSub,
	userContext interface{}) *Reactor {
	LunoPairInformation, _ := userContext.(*common.PairInformation)
	return &Reactor{
		BaseConnectionReactor: common2.NewBaseConnectionReactor(
			logger, cancelCtx, cancelFunc, connectionCancelFunc, userContext),
		messageRouter:       messageRouter.NewMessageRouter(),
		APIKeyID:            APIKeyID,
		APIKeySecret:        APIKeySecret,
		FullMarketOrderBook: fullMarketData.NewFullMarketOrderBook(),
		PubSub:              PubSub,
		LunaPairInformation: LunoPairInformation,
	}
}
