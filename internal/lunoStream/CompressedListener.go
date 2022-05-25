package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/netListener"

	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func CompressedListener(
	ConsumerCounter *netDial.CanDialDefaultImpl,
	maxConnections int, url string, pairInformation ...*common.PairInformation) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	cfrName := "CompressedListenerConnection.CFR"

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
									return listener.NewConnectionReactorFactory(
										cfrName,
										params.PubSub,
										func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
											return stream.Marshall(data)
										},
										ConsumerCounter), nil
								}
							},
						}),
				)
				return netListener.NewNetListenAppNoCrfName(
					fxOptions,
					CompressedListenerConnection,
					url,
					impl.TransportFactoryCompressedTlsName,
					netListener.UserContextValue(pairInformation),
					netListener.MaxConnectionsSetting(maxConnections))(params.NetAppFuncInParams)
			},
		},
	)
}
