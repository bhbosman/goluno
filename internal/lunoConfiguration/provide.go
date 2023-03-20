package lunoConfiguration

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goConn"
	fxAppManager "github.com/bhbosman/goFxAppManager/service"
	"github.com/bhbosman/goLuno/internal/reactorWS"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"net/url"
)

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(
			func(
				params struct {
					fx.In
				},
			) (func() (ILunoConfigurationData, error), error) {
				return func() (ILunoConfigurationData, error) {
					return newData()
				}, nil
			},
		),
		fx.Provide(
			func(
				params struct {
					fx.In
					PubSub                     *pubsub.PubSub  `name:"Application"`
					ApplicationContext         context.Context `name:"Application"`
					OnData                     func() (ILunoConfigurationData, error)
					Lifecycle                  fx.Lifecycle
					Logger                     *zap.Logger
					UniqueReferenceService     interfaces.IUniqueReferenceService
					UniqueSessionNumber        interfaces.IUniqueSessionNumber
					GoFunctionCounter          GoFunctionCounter.IService
					InstrumentReferenceService instrumentReference.IInstrumentReferenceService
				},
			) (ILunoConfigurationService, error) {
				serviceInstance, err := newService(
					params.ApplicationContext,
					params.OnData,
					params.Logger,
					params.PubSub,
					params.GoFunctionCounter,
					params.InstrumentReferenceService,
				)
				if err != nil {
					return nil, err
				}
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: serviceInstance.OnStart,
						OnStop:  serviceInstance.OnStop,
					},
				)
				return serviceInstance, nil
			},
		),
		fx.Invoke(
			func(
				params struct {
					fx.In
					Logger                     *zap.Logger
					LunoConfigurationService   ILunoConfigurationService
					ApplicationContext         context.Context `name:"Application"`
					Lifecycle                  fx.Lifecycle
					FxManagerService           fxAppManager.IFxManagerService
					NetMultiDialer             goCommsMultiDialer.INetMultiDialerService
					PubSub                     *pubsub.PubSub `name:"Application"`
					LunoAPIKeyID               string         `name:"LunoAPIKeyID"`
					LunoAPIKeySecret           string         `name:"LunoAPIKeySecret"`
					FullMarketDataHelper       fullMarketDataHelper.IFullMarketDataHelper
					FmdService                 fullMarketDataManagerService.IFmdManagerService
					InstrumentReferenceService instrumentReference.IInstrumentReferenceService
				},
			) error {
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							allLunoConfiguration, err := params.InstrumentReferenceService.GetLunoProviders() //LunoConfigurationService.GetAll()
							if err != nil {
								return err
							}
							f := func(
								name string,
								referenceData instrumentReference.LunoReferenceData) func() (messages.IApp, goConn.ICancellationContext, error) {

								return func() (messages.IApp, goConn.ICancellationContext, error) {

									namedLogger := params.Logger.Named(name)
									ctx, cancelFunc := context.WithCancel(params.ApplicationContext)
									cancellationContext, err := goConn.NewCancellationContextNoCloser(
										name,
										cancelFunc,
										ctx,
										namedLogger,
									)
									if err != nil {
										return nil, nil, err
									}
									pairUrl, _ := url.Parse(fmt.Sprintf("wss://ws.luno.com:443/api/1/stream/%v", name))

									dec, err := reactorWS.NewDecorator(
										pairUrl,
										params.Logger,
										params.NetMultiDialer,
										referenceData,
										name,
										params.PubSub,
										params.LunoAPIKeyID,
										params.LunoAPIKeySecret,
										params.FullMarketDataHelper,
										params.FmdService,
										cancellationContext,
									)
									if err != nil {
										return nil, nil, err
									}
									return dec, cancellationContext, nil
								}
							}

							for _, configuration := range allLunoConfiguration {
								name := configuration.Name
								err = multierr.Append(
									err,
									params.FxManagerService.Add(
										name,
										f(configuration.Name, configuration),
									),
								)
							}
							return nil
						},
						OnStop: nil,
					},
				)
				return nil
			},
		),
	)
}
