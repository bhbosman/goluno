package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
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

	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					CompressedListenerConnection,
					url,
					impl.TransportFactoryCompressedTlsName,
					netListener.UserContextValue(pairInformation),
					netListener.MaxConnectionsSetting(maxConnections),
					netListener.FxOption(
						fx.Provide(
							fx.Annotated{
								Target: func(pubSub *pubsub.PubSub) (intf.IConnectionReactorFactory, error) {
									return listener.NewConnectionReactorFactory(
										CompressedListenerConnection,
										pubSub,
										func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
											return stream.Marshall(data)
										},
										ConsumerCounter), nil
								}}))),
			}),
	)
}
