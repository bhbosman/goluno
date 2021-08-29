package lunoStream

import (
	"encoding/json"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func TextListener(
	pubSub *pubsub.PubSub,
	ConsumerCounter *netDial.CanDialDefaultImpl,
	maxConnections int, url string, pairInformation ...*common.PairInformation) fx.Option {
	const TextListenerConnection = "TextListenerConnection"
	cfr := listener.NewConnectionReactorFactory(
		TextListenerConnection,
		pubSub,
		func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
			bytes, err := json.MarshalIndent(m, "", "\t")
			if err != nil {
				return nil, err
			}
			return gomessageblock.NewReaderWriterBlock(bytes), nil
		},
		ConsumerCounter)

	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					TextListenerConnection,
					url,
					impl.TransportFactoryEmptyName,
					impl.GetNamedStack(impl.TransportFactoryEmptyName),
					cfr,
					netListener.UserContextValue(pairInformation),
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}
