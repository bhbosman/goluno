package reactorWS

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/goCommsStacks/websocket"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/url"
)

type decorator struct {
	netMultiDialer                goCommsMultiDialer.INetMultiDialerService
	name                          string
	pubSub                        *pubsub.PubSub
	lunoAPIKeyID                  string
	lunoAPIKeySecret              string
	logger                        *zap.Logger
	fullMarketDataHelper          fullMarketDataHelper.IFullMarketDataHelper
	fmdService                    fullMarketDataManagerService.IFmdManagerService
	referenceData                 instrumentReference.LunoReferenceData
	decoratorCancellationContext  gocommon.ICancellationContext
	connectionCancellationContext gocommon.ICancellationContext
	connectionUrl                 *url.URL
}

func (self *decorator) Start(context.Context) error {
	if self.decoratorCancellationContext.Err() != nil {
		return self.decoratorCancellationContext.Err()
	}
	go func() {
		_ = self.internalStart()
	}()
	return nil
}

func (self *decorator) Stop(context.Context) error {
	self.decoratorCancellationContext.Cancel("ASDADSA")
	return self.decoratorCancellationContext.Err()
}

func (self *decorator) Err() error {
	return self.decoratorCancellationContext.Err()
}

func (self *decorator) internalStart() error {
	dialApp, dialAppCancelFunc, connectionId, err := self.netMultiDialer.Dial(
		false,
		nil,
		self.connectionUrl,
		self.reconnect,
		self.decoratorCancellationContext,
		fmt.Sprintf("Luno.%v", self.name),
		fmt.Sprintf("Luno.%v", self.name),
		Provide(),
		goCommsDefinitions.ProvideTransportFactoryForWebSocketName(
			topStack.Provide(),
			websocket.Provide(),
			bottom.Provide(),
		),
		goCommsDefinitions.ProvideStringContext("Pair", self.name),
		goCommsDefinitions.ProvideStringContext("LunoAPIKeyID", self.lunoAPIKeyID),
		goCommsDefinitions.ProvideStringContext("LunoAPIKeySecret", self.lunoAPIKeySecret),
		PubSub.ProvidePubSubInstance("Application", self.pubSub),
		fx.Supply(self.referenceData),
		fx.Provide(
			fx.Annotated{
				Target: func() (fullMarketDataHelper.IFullMarketDataHelper, fullMarketDataManagerService.IFmdManagerService) {
					return self.fullMarketDataHelper, self.fmdService
				},
			},
		),
	)
	if err != nil {
		return err
	}
	err = dialApp.Start(context.Background())
	if err != nil {
		self.logger.Error("Error in start", zap.Error(err))
	}
	self.connectionCancellationContext = dialAppCancelFunc

	return gocommon.RegisterConnectionShutdown(
		connectionId,
		func(
			connectionApp gocommon.IApp,
			logger *zap.Logger,
			connectionCancellationContext gocommon.ICancellationContext,
		) func() {
			return func() {
				errInGoRoutine := connectionApp.Stop(context.Background())
				connectionCancellationContext.Cancel("asadasdas")
				if errInGoRoutine != nil {
					logger.Error(
						"Stopping error. not really a problem. informational",
						zap.Error(errInGoRoutine))
				}
			}
		}(dialApp, self.logger, dialAppCancelFunc),
		self.decoratorCancellationContext,
		dialAppCancelFunc,
	)
}

func (self *decorator) internalStop() error {
	self.connectionCancellationContext.Cancel("adsdasdas")
	self.connectionCancellationContext = nil
	return self.decoratorCancellationContext.Err()
}

func (self *decorator) reconnect() {
	if self.decoratorCancellationContext.Err() != nil {
		return
	}
	go func() {
		err := self.internalStop()
		if err != nil {
			return
		}
		err = self.internalStart()
		if err != nil {
			return
		}
	}()
}

func NewDecorator(
	connectionUrl *url.URL,
	Logger *zap.Logger,
	NetMultiDialer goCommsMultiDialer.INetMultiDialerService,
	referenceData instrumentReference.LunoReferenceData,
	name string,
	pubSub *pubsub.PubSub,
	LunoAPIKeyID string,
	LunoAPIKeySecret string,
	FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper,
	FmdService fullMarketDataManagerService.IFmdManagerService,
	decoratorCancellationContext gocommon.ICancellationContext,
) (gocommon.IApp, error) {
	return &decorator{
		connectionUrl:                connectionUrl,
		logger:                       Logger,
		netMultiDialer:               NetMultiDialer,
		name:                         name,
		pubSub:                       pubSub,
		lunoAPIKeyID:                 LunoAPIKeyID,
		lunoAPIKeySecret:             LunoAPIKeySecret,
		fullMarketDataHelper:         FullMarketDataHelper,
		fmdService:                   FmdService,
		referenceData:                referenceData,
		decoratorCancellationContext: decoratorCancellationContext,
	}, nil
}
