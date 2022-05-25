package lunoStream

import (
	"encoding/json"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func TextListener(
	ConsumerCounter *netDial.CanDialDefaultImpl,
	maxConnections int, url string, pairInformation ...*common.PairInformation) fx.Option {
	const TextListenerConnection = "TextListenerConnection"
	cfrName := "TextListenerConnection"

	return fx.Provide(
		fx.Annotated{
			Group: "Apps",
			Target: func(params struct {
				fx.In
				PubSub             *pubsub.PubSub `name:"Application"`
				NetAppFuncInParams impl.NetAppFuncInParams
			}) messages.CreateAppCallback {
				fxOptions := fx.Options(
					fx.Provide(fx.Annotated{Name: "Application", Target: func() *pubsub.PubSub { return params.PubSub }}),
					fx.Provide(
						fx.Annotated{
							Target: func(params struct {
								fx.In
								PubSub *pubsub.PubSub `name:"Application"`
							}) intf.ConnectionReactorFactoryCallback {
								return func() (intf.IConnectionReactorFactory, error) {
									cfr := listener.NewConnectionReactorFactory(
										cfrName,
										params.PubSub,
										func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
											bytes, err := json.MarshalIndent(m, "", "\t")
											if err != nil {
												return nil, err
											}
											return gomessageblock.NewReaderWriterBlock(bytes), nil
										},
										ConsumerCounter)
									return cfr, nil
								}
							},
						}),
				)
				return netListener.NewNetListenAppNoCrfName(
					fxOptions,
					TextListenerConnection,
					url,
					impl.TransportFactoryEmptyName,
					netListener.UserContextValue(pairInformation),
					netListener.MaxConnectionsSetting(2000))(params.NetAppFuncInParams)
			},
		})

}
