package reactorWS

import (
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/goCommsStacks/websocket"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"io"
	"net/url"
)

type decorator struct {
	stoppedCalled        bool
	NetMultiDialer       goCommsMultiDialer.INetMultiDialerService
	name                 string
	pubSub               *pubsub.PubSub
	LunoAPIKeyID         string
	LunoAPIKeySecret     string
	dialApp              messages.IApp
	dialAppCancelFunc    goCommsDefinitions.ICancellationContext
	Logger               *zap.Logger
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper
	FmdService           fullMarketDataManagerService.IFmdManagerService
	referenceData        instrumentReference.LunoReferenceData
}

func (self *decorator) Cancel() {
	if self.dialAppCancelFunc != nil {
		self.dialAppCancelFunc.Cancel()
	}
}

func (self *decorator) Start(ctx context.Context) error {
	if !self.stoppedCalled {
		go func() {
			_ = self.internalStart(ctx)
		}()
		return nil
	}
	return io.EOF
}

func (self *decorator) Stop(ctx context.Context) error {
	if !self.stoppedCalled {
		self.stoppedCalled = true
		return self.internalStop(ctx)
	}
	return io.EOF
}

func (self *decorator) Err() error {
	if self.dialApp != nil {
		return self.dialApp.Err()
	}
	return nil
}

func (self *decorator) internalStart(ctx context.Context) error {
	pairUrl, _ := url.Parse(fmt.Sprintf("wss://ws.luno.com:443/api/1/stream/%v", self.name))
	var err error
	var connectionId string
	self.dialApp, self.dialAppCancelFunc, connectionId, err = self.NetMultiDialer.Dial(
		false,
		nil,
		pairUrl,
		self.reconnect,
		self.dialAppCancelFunc,
		fmt.Sprintf("Luno.%v", self.name),
		fmt.Sprintf("Luno.%v", self.name),
		Provide(),
		goCommsDefinitions.ProvideTransportFactoryForWebSocketName(
			topStack.ProvideTopStack(),
			websocket.ProvideWebsocketStacks(),
			bottom.Provide(),
		),
		goCommsDefinitions.ProvideStringContext("Pair", self.name),
		goCommsDefinitions.ProvideStringContext("LunoAPIKeyID", self.LunoAPIKeyID),
		goCommsDefinitions.ProvideStringContext("LunoAPIKeySecret", self.LunoAPIKeySecret),
		PubSub.ProvidePubSubInstance("Application", self.pubSub),
		fx.Supply(self.referenceData),
		fx.Provide(
			fx.Annotated{
				Target: func() (fullMarketDataHelper.IFullMarketDataHelper, fullMarketDataManagerService.IFmdManagerService) {
					return self.FullMarketDataHelper, self.FmdService
				},
			},
		),
	)
	if err != nil {
		return err
	}
	err = self.dialApp.Start(context.Background())
	if err != nil {
		self.Logger.Error("Error in start", zap.Error(err))
	}

	err = self.dialAppCancelFunc.Add(
		connectionId,
		func() func() {
			b := false
			return func() {
				if !b {
					b = true
					stopErr := self.dialApp.Stop(context.Background())
					if stopErr != nil {
						self.Logger.Error(
							"Stopping error. not really a problem. informational",
							zap.Error(stopErr))
					}
					_ = self.dialAppCancelFunc.Remove(connectionId)
				}
			}
		}(),
	)
	return nil
}

func (self *decorator) internalStop(ctx context.Context) error {
	if self.dialAppCancelFunc != nil {
		self.dialAppCancelFunc.Cancel()
	}
	return nil
}

func (self *decorator) reconnect() {
	go func() {
		if !self.stoppedCalled {
			err := self.internalStop(context.Background())
			if err != nil {
				return
			}
			err = self.internalStart(context.Background())
			if err != nil {
				return
			}
		}
	}()
}

func NewDecorator(
	Logger *zap.Logger,
	NetMultiDialer goCommsMultiDialer.INetMultiDialerService,
	referenceData instrumentReference.LunoReferenceData,
	name string,
	pubSub *pubsub.PubSub,
	LunoAPIKeyID string,
	LunoAPIKeySecret string,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
	FmdService fullMarketDataManagerService.IFmdManagerService,
) *decorator {
	return &decorator{
		Logger:               Logger,
		NetMultiDialer:       NetMultiDialer,
		name:                 name,
		pubSub:               pubSub,
		LunoAPIKeyID:         LunoAPIKeyID,
		LunoAPIKeySecret:     LunoAPIKeySecret,
		FullMarketDataHelper: FullMarketDataHelper,
		FmdService:           FmdService,
		referenceData:        referenceData,
	}
}
