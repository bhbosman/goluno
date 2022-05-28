package lunoWS

import (
	"fmt"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func ProvideDialers(options ...DialersApply) fx.Option {
	settings := &lunoStreamDialersSettings{}
	for _, option := range options {
		if option == nil {
			continue
		}
		option.apply(settings)
	}
	var opt []fx.Option
	for _, optionEnum := range settings.pairs {
		if optionEnum == nil {
			continue
		}

		option := optionEnum

		crfName := fmt.Sprintf("LunoStream.%v.CRF", option.Pair)

		opt = append(
			opt,
			fx.Provide(
				fx.Annotated{
					Group: "Apps",
					Target: func(params struct {
						fx.In
						PubSub           *pubsub.PubSub `name:"Application"`
						LunoAPIKeyID     string         `name:"LunoAPIKeyID"`
						LunoAPIKeySecret string         `name:"LunoAPIKeySecret"`
						AppFuncInParams  impl.NetAppFuncInParams
					}) messages.CreateAppCallback {
						fxOptions := fx.Options(
							fx.Provide(fx.Annotated{Name: "LunoAPIKeyID", Target: model.CreateStringContext(params.LunoAPIKeyID)}),
							fx.Provide(fx.Annotated{Name: "LunoAPIKeySecret", Target: model.CreateStringContext(params.LunoAPIKeySecret)}),
							fx.Provide(fx.Annotated{Name: "Application", Target: func() *pubsub.PubSub { return params.PubSub }}),
							fx.Provide(
								fx.Annotated{
									Target: func(params struct {
										fx.In
										PubSub           *pubsub.PubSub `name:"Application"`
										LunoAPIKeyID     string         `name:"LunoAPIKeyID"`
										LunoAPIKeySecret string         `name:"LunoAPIKeySecret"`
									}) intf.ConnectionReactorFactoryCallback {
										return func() (intf.IConnectionReactorFactory, error) {
											return NewConnectionReactorFactory(
												crfName,
												params.LunoAPIKeyID,
												params.LunoAPIKeySecret,
												params.PubSub), nil

										}
									},
								},
							),
						)

						return netDial.NewNetDialAppNoCrfName(
							fxOptions,
							fmt.Sprintf("LunoStream[%v]", option.Pair),
							fmt.Sprintf("wss://ws.luno.com:443/api/1/stream/%v", option.Pair),
							impl.WebSocketName,
							netDial.MaxConnectionsSetting(settings.maxConnections),
							netDial.UserContextValue(option),
							netDial.CanDial(settings.canDial...))(params.AppFuncInParams)
					},
				}))
	}
	return fx.Options(opt...)
}
