package lunoConfiguration

import (
	"context"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsMultiDialer"
	fxAppManager "github.com/bhbosman/goFxAppManager/service"
	"github.com/bhbosman/goLuno/internal/reactorWS"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
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
								referenceData instrumentReference.LunoReferenceData) func() (messages.IApp, context.CancelFunc, error) {
								return func() (messages.IApp, context.CancelFunc, error) {
									dec := reactorWS.NewDecorator(
										params.Logger,
										params.NetMultiDialer,
										referenceData,
										name,
										params.PubSub,
										params.LunoAPIKeyID,
										params.LunoAPIKeySecret,
										params.FullMarketDataHelper,
										params.FmdService,
									)
									return dec, dec.Cancel, nil
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
