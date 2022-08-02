package lunoConfiguration

import (
	"context"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type service struct {
	parentContext              context.Context
	ctx                        context.Context
	cancelFunc                 context.CancelFunc
	cmdChannel                 chan interface{}
	onData                     func() (ILunoConfigurationData, error)
	Logger                     *zap.Logger
	state                      IFxService.State
	pubSub                     *pubsub.PubSub
	goFunctionCounter          GoFunctionCounter.IService
	subscribeChannel           *pubsub.NextFuncSubscription
	InstrumentReferenceService instrumentReference.IInstrumentReferenceService
}

func (self *service) Send(message interface{}) error {
	send, err := CallILunoConfigurationSend(self.ctx, self.cmdChannel, false, message)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *service) OnStart(ctx context.Context) error {
	err := self.start(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *service) OnStop(ctx context.Context) error {
	err := self.shutdown(ctx)
	close(self.cmdChannel)
	self.state = IFxService.Stopped
	return err
}

func (self *service) shutdown(_ context.Context) error {
	self.cancelFunc()
	return pubSub.Unsubscribe("", self.pubSub, self.goFunctionCounter, self.subscribeChannel)
}

func (self *service) start(_ context.Context) error {
	instanceData, err := self.onData()
	if err != nil {
		return err
	}

	err = self.goFunctionCounter.GoRun(
		"Luno Configuration Service",
		func() {
			self.goStart(instanceData)
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (self *service) goStart(instanceData ILunoConfigurationData) {
	self.subscribeChannel = pubsub.NewNextFuncSubscription(goCommsDefinitions.CreateNextFunc(self.cmdChannel))
	self.pubSub.AddSub(self.subscribeChannel, self.ServiceName())

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		instanceData,
		[]ChannelHandler.ChannelHandler{
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ILunoConfiguration); ok {
						return ChannelEventsForILunoConfiguration(unk, message)
					}
					return false, nil
				},
			},
			{
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return true, unk.Send(message)
					}
					return false, nil
				},
			},
		},
		func() int {
			return len(self.cmdChannel) + self.subscribeChannel.Count()
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
	)
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := instanceData.ShutDown()
			if err != nil {
				self.Logger.Error(
					"error on done",
					zap.Error(err))
			}
			break loop
		case event, ok := <-self.cmdChannel:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break loop
			}
		}
	}
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self service) ServiceName() string {
	return "LunoConfiguration"
}

func newService(
	parentContext context.Context,
	onData func() (ILunoConfigurationData, error),
	logger *zap.Logger,
	pubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
	InstrumentReferenceService instrumentReference.IInstrumentReferenceService,
) (ILunoConfigurationService, error) {
	localCtx, localCancelFunc := context.WithCancel(parentContext)
	return &service{
			parentContext:              parentContext,
			ctx:                        localCtx,
			cancelFunc:                 localCancelFunc,
			cmdChannel:                 make(chan interface{}, 32),
			onData:                     onData,
			Logger:                     logger,
			pubSub:                     pubSub,
			goFunctionCounter:          goFunctionCounter,
			InstrumentReferenceService: InstrumentReferenceService,
		},
		nil
}
