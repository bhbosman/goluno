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

func CompressedListener(maxConnections int, url string, pairInformation ...*common.PairInformation) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: impl.ConnectionReactorFactoryConst,
				Target: func(params struct {
					fx.In
					PubSub          *pubsub.PubSub `name:"Application"`
					ConsumerCounter *netDial.CanDialDefaultImpl
				}) (intf.IConnectionReactorFactory, error) {
					return listener.NewConnectionReactorFactory(
						CompressedListenerConnection,
						params.PubSub,
						func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
							return stream.Marshall(data)
						},
						params.ConsumerCounter), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					CompressedListenerConnection,
					url,
					impl.TransportFactoryCompressedName,
					CompressedListenerConnection,
					netListener.UserContextValue(pairInformation),
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}
