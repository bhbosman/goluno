package lunoStream

import (
	"encoding/json"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
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

func TextListener(maxConnections int, url string, pairInformation ...*common.PairInformation) fx.Option {
	const TextListenerConnection = "TextListenerConnection"
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
					impl.TransportFactoryEmptyName,
					TextListenerConnection,
					netListener.UserContextValue(pairInformation),
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}
