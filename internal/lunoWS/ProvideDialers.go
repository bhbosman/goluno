package lunoWS

import (
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/top"
	"github.com/bhbosman/goCommsStacks/websocket"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/common"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func ProvideDialers(
	options ...DialersApply,
) fx.Option {
	settings := &lunoStreamDialersSettings{}
	for _, option := range options {
		if option == nil {
			continue
		}
		option.apply(settings)
	}
	var opt []fx.Option
	for _, lunoPair := range settings.pairs {
		if lunoPair == nil {
			continue
		}
		option := lunoPair
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
						AppFuncInParams  common.NetAppFuncInParams
					}) messages.CreateAppCallback {
						f := goCommsNetDialer.NewNetDialApp(
							fmt.Sprintf("LunoStream.%v.ConnectionManager", option.Pair),
							option.ServiceIdentifier,
							option.ServiceDependentOn,
							fmt.Sprintf("LunoStream.%v.Connection", option.Pair),
							option.UseSocks5,
							option.SocksUrl,
							option.PairUrl, // fmt.Sprintf("wss://ws.luno.com:443/api/1/stream/%v", option.Pair),
							goCommsDefinitions.WebSocketName,
							common.MaxConnectionsSetting(settings.maxConnections),
							goCommsNetDialer.UserContextValue(option),
							goCommsNetDialer.CanDial(settings.canDial...),

							// TODO: work out IConnectionReactorFactory
							//common.NewOverrideCreateConnectionReactor(
							//	fx.Provide(
							//		fx.Annotated{
							//			Target: func(
							//				params struct {
							//					fx.In
							//					LunoAPIKeyID         string         `name:"LunoAPIKeyID"`
							//					LunoAPIKeySecret     string         `name:"LunoAPIKeySecret"`
							//					PubSub               *pubsub.PubSub `name:"Application"`
							//					CancelCtx            context.Context
							//					CancelFunc           context.CancelFunc
							//					ConnectionCancelFunc model.ConnectionCancelFunc
							//					Logger               *zap.Logger
							//					Cfr                  intf.IConnectionReactorFactory
							//					ClientContext        interface{} `name:"UserContext"`
							//				},
							//			) (intf.IConnectionReactor, error) {
							//				result := NewConnectionReactor(
							//					params.Logger,
							//					params.CancelCtx,
							//					params.CancelFunc,
							//					params.ConnectionCancelFunc,
							//					params.LunoAPIKeyID,
							//					params.LunoAPIKeySecret,
							//					params.PubSub,
							//					params.ClientContext,
							//				)
							//				return result, nil
							//
							//			},
							//		},
							//	),
							//),
							common.NewConnectionInstanceOptions(
								goCommsDefinitions.ProvideTransportFactoryForWebSocketName(
									top.ProvideTopStack(),
									bottom.Provide(),
									websocket.ProvideWebsocketStacks(),
								),
								goCommsDefinitions.ProvideStringContext("LunoAPIKeyID", params.LunoAPIKeyID),
								goCommsDefinitions.ProvideStringContext("LunoAPIKeySecret", params.LunoAPIKeySecret),
								PubSub.ProvidePubSubInstance("Application", params.PubSub),
								ProvideConnectionReactorFactory(),
							),
						)

						return f(
							params.AppFuncInParams)
					},
				},
			),
		)
	}
	return fx.Options(opt...)
}
