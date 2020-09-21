package lunoStream

import (
	"encoding/json"
	"github.com/bhbosman/goLuno/internal/ConsumerCounter"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/netListener"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func TextListener(url string, pairInformation ...*common.PairInformation) fx.Option {
	const TextListenerConnection = "TextListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: commsImpl.ConnectionReactorFactoryConst,
				Target: func(params struct {
					fx.In
					PubSub          *pubsub.PubSub `name:"Application"`
					ConsumerCounter *ConsumerCounter.ConsumerCounter
				}) (commsImpl.IConnectionReactorFactory, error) {
					return listener.NewConnectionReactorFactory(
						TextListenerConnection,
						params.PubSub,
						func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
							bytes, err := json.MarshalIndent(m, "", "\t")
							if err != nil {
								return nil, err
							}
							return gomessageblock.NewReaderWriterBlock(bytes), nil
						},
						params.ConsumerCounter), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					TextListenerConnection,
					url,
					commsImpl.TransportFactoryEmptyName,
					TextListenerConnection,
					netListener.UserContextValue(pairInformation)),
			}),
	)
}
