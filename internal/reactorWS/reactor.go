package reactorWS

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketData"
	stream2 "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	lunaRawDataFeed "github.com/bhbosman/goMessages/luno/stream"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messages"
	common3 "github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gomessageblock"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"

	"github.com/bhbosman/gocommon/messageRouter"

	"github.com/bhbosman/goCommsStacks/webSocketMessages/wsmsg"
	"github.com/bhbosman/gocommon/stream"
	"github.com/cskr/pubsub"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type reactor struct {
	common2.BaseConnectionReactor
	apiKeyID      string
	apiKeySecret  string
	updateCount   int64
	sequence      int64
	messageRouter messageRouter.IMessageRouter
	fmdService    fullMarketDataManagerService.IFmdManagerService
	fmdHelper     fullMarketDataHelper.IFullMarketDataHelper
	referenceData instrumentReference.LunoReferenceData
}

func (self *reactor) SendMessage(message proto.Message) error {
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
	wsMsg := wsmsg.WebSocketMessage{
		OpCode:  wsmsg.WebSocketMessage_OpText,
		Message: flatten,
	}
	readWriterSize, err := stream.Marshall(&wsMsg)
	if err != nil {
		return err
	}
	self.OnSendToConnection(readWriterSize)
	return nil
}

func (self *reactor) Init(
	params intf.IInitParams,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	_, _, _, err := self.BaseConnectionReactor.Init(params)
	if err != nil {
		return nil, nil, nil, err
	}

	return func(i interface{}) {
			self.messageRouter.Route(i)
		},
		func(err error) {
			self.messageRouter.Route(err)
		},
		func() {
			self.CancelFunc()
		},
		nil
}

func (self *reactor) Close() error {
	_ = self.fmdService.Send(
		&stream2.FullMarketData_RemoveInstrumentInstruction{
			Instrument: self.referenceData.SystemName,
		},
	)

	return self.BaseConnectionReactor.Close()
}

func (self *reactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	if err != nil {
		return err
	}
	self.fmdService.MultiSend(
		&stream2.FullMarketData_Clear{
			Instrument: self.referenceData.SystemName,
		},
		&stream2.FullMarketData_Instrument_InstrumentStatus{
			Instrument: self.referenceData.SystemName,
			Status:     "Connection opened",
		},
	)
	return nil
}

func (self *reactor) HandleReaderWriter(msg *gomessageblock.ReaderWriter) {
	marshal, err := stream.UnMarshal(msg)
	if err != nil {
		return
	}
	self.messageRouter.Route(marshal)
	if err != nil {
		return
	}
}

func (self *reactor) HandleEmptyQueue(_ *messages.EmptyQueue) {
}

func (self *reactor) updateSequence(newSeq int64) error {
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

func (self *reactor) HandleLunaData(msg *lunaRawDataFeed.LunoStreamData) error {
	err := self.updateSequence(msg.Sequence)
	if err != nil {
		self.CancelFunc()
		return err
	}
	switch {
	case msg.Bids != nil || msg.Asks != nil:
		clearMessage := &stream2.FullMarketData_Clear{
			Instrument: self.referenceData.SystemName,
		}
		_ = self.fmdService.Send(
			clearMessage,
		)

		for _, bidOrder := range msg.Bids {
			instruction := &stream2.FullMarketData_AddOrderInstruction{
				Instrument: self.referenceData.SystemName,
				Order: &stream2.FullMarketData_AddOrder{
					Side:   fullMarketData.BuySide,
					Id:     bidOrder.Id,
					Price:  bidOrder.Price,
					Volume: bidOrder.Volume,
				},
			}
			_ = self.fmdService.Send(
				instruction,
			)
		}
		for _, askOrder := range msg.Asks {
			instruction := &stream2.FullMarketData_AddOrderInstruction{
				Instrument: self.referenceData.SystemName,
				Order: &stream2.FullMarketData_AddOrder{
					Side:   fullMarketData.AskSide,
					Id:     askOrder.Id,
					Price:  askOrder.Price,
					Volume: askOrder.Volume,
				},
			}
			_ = self.fmdService.Send(
				instruction,
			)
		}

	case msg.TradeUpdates != nil:
		for _, order := range msg.TradeUpdates {
			reduceVolume := &stream2.FullMarketData_ReduceVolumeInstruction{
				Instrument: self.referenceData.SystemName,
				Id:         order.MakerOrderId,
				Volume:     order.Base,
			}
			_ = self.fmdService.Send(
				reduceVolume,
			)
		}
	case msg.DeleteUpdate != nil:
		deleteOrder := &stream2.FullMarketData_DeleteOrderInstruction{
			Instrument: self.referenceData.SystemName,
			Id:         msg.DeleteUpdate.OrderId,
		}
		_ = self.fmdService.Send(
			deleteOrder,
		)
	case msg.CreateUpdate != nil:
		order := func(createUpdate *lunaRawDataFeed.CreateUpdate) *stream2.FullMarketData_AddOrder {
			if createUpdate.Type == "BID" {
				return &stream2.FullMarketData_AddOrder{
					Side:   stream2.OrderSide_BidOrder,
					Id:     createUpdate.Id,
					Price:  createUpdate.Price,
					Volume: createUpdate.Volume,
				}
			}
			return &stream2.FullMarketData_AddOrder{
				Side:   stream2.OrderSide_AskOrder,
				Id:     createUpdate.Id,
				Price:  createUpdate.Price,
				Volume: createUpdate.Volume,
			}
		}(msg.CreateUpdate)
		_ = self.fmdService.Send(
			&stream2.FullMarketData_AddOrderInstruction{
				Instrument: self.referenceData.SystemName,
				Order:      order,
			},
		)
	}
	return nil
}
func (self *reactor) HandleWebSocketMessageWrapper(msg *wsmsg.WebSocketMessageWrapper) {
	self.HandleWebSocketMessage(msg.Data)
}

type LunoError struct {
	ErrorMessage string
	ErrorCode    string
}

func (self *LunoError) Error() string {
	return fmt.Sprintf("error: %v, error_code: %v", self.Error, self.ErrorCode)
}

func NewLunoError(errorMessage string, errorCode string) *LunoError {
	return &LunoError{
		ErrorMessage: errorMessage,
		ErrorCode:    errorCode,
	}
}

func (self *reactor) HandleWebSocketMessage(msg *wsmsg.WebSocketMessage) {
	switch msg.OpCode {
	case wsmsg.WebSocketMessage_OpText:
		Unmarshaler := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
			AnyResolver:        nil,
		}
		if len(msg.Message) == 2 && string(msg.Message) == "\"\"" {
			// this is a heartbeat as per spec
			// An empty message is a keep alive message.
			return
		}
		lunaData := &lunaRawDataFeed.LunoStreamData{}
		err := Unmarshaler.Unmarshal(bytes.NewBuffer(msg.Message), lunaData)
		if err != nil {
			self.ConnectionCancelFunc("WSLuno, json unmarshall error", true, err)
			return
		}
		if lunaData.Status != "" {
			_ = self.fmdService.Send(
				&stream2.FullMarketData_Instrument_InstrumentStatus{
					Instrument: self.referenceData.SystemName,
					Status:     lunaData.Status,
				},
			)
		}

		if lunaData.ErrorCode == "sessions_limit_exceeded" {
			self.ConnectionCancelFunc(
				"",
				true,
				NewLunoError(lunaData.Error, lunaData.ErrorCode),
			)
			self.CancelFunc()
			return
		}
		self.messageRouter.Route(lunaData)
		return
	case wsmsg.WebSocketMessage_OpStartLoop:
		msg := &lunaRawDataFeed.Credentials{
			ApiKeyId:     self.apiKeyID,
			ApiKeySecret: self.apiKeySecret,
		}
		_ = self.SendMessage(msg)
		return
	default:
		return
	}
}

func newConnectionReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc common3.ConnectionCancelFunc,
	apiKeyID string,
	apiKeySecret string,
	pubSub *pubsub.PubSub,
	referenceData instrumentReference.LunoReferenceData,
	goFunctionCounter GoFunctionCounter.IService,
	uniqueReferenceService interfaces.IUniqueReferenceService,
	fullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
	fmdService fullMarketDataManagerService.IFmdManagerService,
) (intf.IConnectionReactor, error) {
	result := &reactor{
		BaseConnectionReactor: common2.NewBaseConnectionReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			uniqueReferenceService.Next("ConnectionReactor"),
			pubSub,
			goFunctionCounter,
		),
		apiKeyID:      apiKeyID,
		apiKeySecret:  apiKeySecret,
		messageRouter: messageRouter.NewMessageRouter(),
		fmdService:    fmdService,
		fmdHelper:     fullMarketDataHelper,
		referenceData: referenceData,
	}
	_ = result.messageRouter.Add(result.HandleWebSocketMessageWrapper)
	_ = result.messageRouter.Add(result.HandleWebSocketMessage)
	_ = result.messageRouter.Add(result.HandleReaderWriter)
	_ = result.messageRouter.Add(result.HandleLunaData)
	_ = result.messageRouter.Add(result.HandleEmptyQueue)

	return result, nil
}
