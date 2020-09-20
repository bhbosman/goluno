package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/ConsumerCounter"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/netListener"

	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)



const CompressedListenerConnection = "CompressedListenerConnection"
func CompressedListener(url string, pairInformation ...*common.PairInformation) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: commsImpl.ConnectionReactorFactoryConst,
				Target: func(params struct {
					fx.In
					PubSub *pubsub.PubSub `name:"Application"`
					ConsumerCounter *ConsumerCounter.ConsumerCounter
				}) (commsImpl.IConnectionReactorFactory, error) {
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
					commsImpl.TransportFactoryCompressedName,
					CompressedListenerConnection,
					netListener.UserContextValue(pairInformation)),
			}),
	)
}
