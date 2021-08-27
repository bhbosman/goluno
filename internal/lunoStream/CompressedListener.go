package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/netListener"

	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func CompressedListener(
	pubSub *pubsub.PubSub,
	ConsumerCounter *netDial.CanDialDefaultImpl,
	maxConnections int, url string, pairInformation ...*common.PairInformation) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	cfr := listener.NewConnectionReactorFactory(
		CompressedListenerConnection,
		pubSub,
		func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
			return stream.Marshall(data)
		},
		ConsumerCounter)
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					CompressedListenerConnection,
					url,
					impl.CreateCompressedStack,
					//impl.CreateCompressedTlsStack,

					cfr,
					netListener.UserContextValue(pairInformation),
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}
