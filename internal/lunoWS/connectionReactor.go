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
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/gomessageblock"

	"github.com/bhbosman/gocommon/messageRouter"

	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/stacks/websocket/wsmsg"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/reactivex/rxgo/v2"
	log2 "log"
	"net"
	"net/url"
)

type ConnectionReactor struct {
	impl.BaseConnectionReactor
	messageRouter        *messageRouter.MessageRouter
	APIKeyID             string
	APIKeySecret         string
	FullMarketOrderBook  *fullMarketData.FullMarketOrderBook
	PubSub               *pubsub.PubSub
	LunoPairInformation  *common.PairInformation
	republishChannelName string
	publishChannelName   string
	UpdateCount          int64
	sequence             int64
}

func (self *ConnectionReactor) SendMessage(message proto.Message) error {
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

func (self *ConnectionReactor) Init(
	conn net.Conn,
	url *url.URL,
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	toConnectionFunc goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc) (rxgo.NextExternalFunc, error) {
	_, _ = self.BaseConnectionReactor.Init(conn, url, connectionId, connectionManager, toConnectionFunc, toConnectionReactor)
	self.Logger.NameChange(fmt.Sprintf("Luno: %v", self.LunoPairInformation.Pair))
	self.ConnectionManager.NameConnection(self.ConnectionId, fmt.Sprintf("Luno: %v", self.LunoPairInformation.Pair))
	_ = self.messageRouter.Add(self.HandleWebSocketMessageWrapper)
	_ = self.messageRouter.Add(self.HandleReaderWriter)
	_ = self.messageRouter.Add(self.HandleLunaData)
	_ = self.messageRouter.Add(self.HandleEmptyQueue)
	_ = self.messageRouter.Add(self.HandlePublishMessage)

	self.republishChannelName = common.RepublishName(self.LunoPairInformation.Pair)
	self.publishChannelName = common.PublishName(self.LunoPairInformation.Pair)

	republishChannel := self.PubSub.Sub(self.republishChannelName)
	go func(ch chan interface{}, topics ...string) {
		defer self.PubSub.Unsub(ch, topics...)
		for {
			select {
			case <-self.CancelCtx.Done():
				return
			case <-ch:
				_ = self.ToReactor(false, &internal.PublishMessage{})
			}
		}
	}(republishChannel, self.republishChannelName)

	return self.doNext, nil
}

func (self *ConnectionReactor) Close() error {
	top5 := &marketDataStream.PublishTop5{
		Instrument: self.LunoPairInformation.Pair,
	}
	self.PubSub.Pub(top5, self.publishChannelName)

	err := self.BaseConnectionReactor.Close()
	return err
}

func (self *ConnectionReactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	return err
}

func (self *ConnectionReactor) doNext(external bool, i interface{}) {
	_, err := self.messageRouter.Route(i)
	if err != nil {
		return
	}
}

func (self *ConnectionReactor) HandleReaderWriter(msg *gomessageblock.ReaderWriter) error {
	marshal, err := stream.UnMarshal(msg, self.CancelCtx, self.CancelFunc, self.ToReactor, self.ToConnection)
	if err != nil {
		return err
	}
	_, err = self.messageRouter.Route(marshal)
	return err
}

func (self *ConnectionReactor) HandlePublishMessage(msg *internal.PublishMessage) error {
	return self.publishData(true)
}

func (self *ConnectionReactor) HandleEmptyQueue(msg *rxgo.EmptyQueue) error {
	self.Logger.LogWithLevel(0, func(logger *log2.Logger) {
	})
	return self.publishData(false)
}
func (self *ConnectionReactor) publishData(forcePublish bool) error {
	thereWasAChange := forcePublish
	maxDepth := 2
	var bids []*marketDataStream.Point
	if highBidNode := self.FullMarketOrderBook.OrderSide[fullMarketData.BuySide].Right(); highBidNode != nil {
		count := 0
		for node := highBidNode; node != nil && count < maxDepth; node = node.Prev() {
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
		for node := lowAskNode; node != nil && count < maxDepth; node = node.Next() {
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
			Instrument:         self.LunoPairInformation.Pair,
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

func (self *ConnectionReactor) HandleLunaData(msg *lunaRawDataFeed.LunoStreamData) error {
	self.FullMarketOrderBook.SetTimeStamp(msg.Timestamp)
	self.FullMarketOrderBook.UpdateMessageReceivedCount()
	switch {
	case msg.Bids != nil || msg.Asks != nil:
		self.sequence = msg.Sequence
		self.FullMarketOrderBook.Clear()
		for _, order := range msg.Bids {
			self.FullMarketOrderBook.AddOrder(fullMarketData.BuySide, order)
		}
		for _, order := range msg.Asks {
			self.FullMarketOrderBook.AddOrder(fullMarketData.AskSide, order)
		}
	case msg.TradeUpdates != nil:
		if self.sequence+1 != msg.Sequence {
			self.CancelFunc()
		}
		self.sequence = msg.Sequence
		for _, order := range msg.TradeUpdates {
			self.FullMarketOrderBook.TradeUpdate(order)
		}
	case msg.DeleteUpdate != nil:
		if self.sequence+1 != msg.Sequence {
			self.CancelFunc()
		}
		self.sequence = msg.Sequence
		self.FullMarketOrderBook.DeleteUpdate(msg.DeleteUpdate)

	case msg.CreateUpdate != nil:
		if self.sequence+1 != msg.Sequence {
			self.CancelFunc()
		}
		self.sequence = msg.Sequence
		self.FullMarketOrderBook.CreateUpdate(msg.CreateUpdate)
	}
	return nil
}

func (self *ConnectionReactor) HandleWebSocketMessageWrapper(msg *wsmsg.WebSocketMessageWrapper) error {
	switch msg.Data.OpCode {
	case wsmsg.WebSocketMessage_OpText:
		unMarshaler := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
			AnyResolver:        nil,
		}

		lunaData := &lunaRawDataFeed.LunoStreamData{}

		err := unMarshaler.Unmarshal(bytes.NewBuffer(msg.Data.Message), lunaData)
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
		self.SendMessage(msg)
		return nil
	default:
		return nil
	}
}

func NewConnectionReactor(
	logger *gologging.SubSystemLogger,
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	APIKeyID string,
	APIKeySecret string,
	PubSub *pubsub.PubSub,
	userContext interface{}) *ConnectionReactor {
	LunoPairInformation, _ := userContext.(*common.PairInformation)
	return &ConnectionReactor{
		BaseConnectionReactor: impl.NewBaseConnectionReactor(
			logger, name, cancelCtx, cancelFunc, userContext),
		messageRouter:       messageRouter.NewMessageRouter(),
		APIKeyID:            APIKeyID,
		APIKeySecret:        APIKeySecret,
		FullMarketOrderBook: fullMarketData.NewFullMarketOrderBook(),
		PubSub:              PubSub,
		LunoPairInformation: LunoPairInformation,
	}
}
