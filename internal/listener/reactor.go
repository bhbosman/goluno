package listener

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketData/stream"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/interfaces"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

type reactor struct {
	common.BaseConnectionReactor
	messageRouter            messageRouter.IMessageRouter
	UniqueReferenceService   interfaces.IUniqueReferenceService
	fmdServiceHelper         fullMarketDataHelper.IFullMarketDataHelper
	fmdService               fullMarketDataManagerService.IFmdManagerService
	registeredFmdInstruments map[string]bool
}

func (self *reactor) Init(params intf.IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	_, _, _, err := self.BaseConnectionReactor.Init(params)
	if err != nil {
		return nil, nil, nil, err
	}

	return func(i interface{}) {
			self.doNext(false, i)
		},
		func(err error) {
			self.doNext(false, err)
		},
		func() {

		},
		nil
}

func (self *reactor) doNext(_ bool, i interface{}) {
	self.messageRouter.Route(i)
}

func (self *reactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	if err != nil {
		return err
	}
	return nil
}

func (self *reactor) Close() error {
	return self.BaseConnectionReactor.Close()
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_Instrument_RegisterWrapper(message *stream.FullMarketData_Instrument_RegisterWrapper) {
	key := self.fmdServiceHelper.InstrumentChannelName(message.Data.Instrument)
	self.PubSub.AddSub(self.OnSendToConnectionPubSubBag, key)
	message.SetNext(self.OnSendToConnectionPubSubBag)
	_ = self.fmdService.Send(message)
	self.registeredFmdInstruments[message.Data.Instrument] = true
}

func (self *reactor) handlePublishRxHandlerCounters(message *model.PublishRxHandlerCounters) {
	for s := range self.registeredFmdInstruments {
		message.AddMapData(fmt.Sprintf("FMD: %v", s), "")
	}
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_Instrument_UnregisterWrapper(message *stream.FullMarketData_Instrument_UnregisterWrapper) {
	key := self.fmdServiceHelper.InstrumentChannelName(message.Data.Instrument)
	self.PubSub.Unsub(self.OnSendToConnectionPubSubBag, key)
	message.SetNext(self.OnSendToConnectionPubSubBag)
	_ = self.fmdService.Send(message)
	delete(self.registeredFmdInstruments, message.Data.Instrument)
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_InstrumentList_SubscribeWrapper(*stream.FullMarketData_InstrumentList_SubscribeWrapper) {
	self.PubSub.AddSub(
		self.OnSendToConnectionPubSubBag,
		self.fmdServiceHelper.InstrumentListChannelName(),
	)
}

//goland:noinspection GoSnakeCaseUsage
func (self *reactor) handleFullMarketData_InstrumentList_RequestWrapper(message *stream.FullMarketData_InstrumentList_RequestWrapper) {
	message.SetNext(self.OnSendToConnectionPubSubBag)
	_ = self.fmdService.Send(message)
}

func newConnectionReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	PubSub *pubsub.PubSub,
	GoFunctionCounter GoFunctionCounter.IService,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
	FmdService fullMarketDataManagerService.IFmdManagerService,
) (intf.IConnectionReactor, error) {
	result := &reactor{
		BaseConnectionReactor: common.NewBaseConnectionReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			UniqueReferenceService.Next("ConnectionReactor"),
			PubSub,
			GoFunctionCounter,
		),
		messageRouter:            messageRouter.NewMessageRouter(),
		UniqueReferenceService:   UniqueReferenceService,
		fmdService:               FmdService,
		fmdServiceHelper:         FullMarketDataHelper,
		registeredFmdInstruments: make(map[string]bool),
	}
	_ = result.messageRouter.Add(result.handleFullMarketData_InstrumentList_SubscribeWrapper)
	_ = result.messageRouter.Add(result.handleFullMarketData_InstrumentList_RequestWrapper)
	//
	_ = result.messageRouter.Add(result.handleFullMarketData_Instrument_RegisterWrapper)
	_ = result.messageRouter.Add(result.handleFullMarketData_Instrument_UnregisterWrapper)
	_ = result.messageRouter.Add(result.handlePublishRxHandlerCounters)

	return result, nil
}
