package provide

import (
	"encoding/json"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/multiBlock"
	"github.com/bhbosman/gocommon/stream"
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
					PubSub *pubsub.PubSub `name:"Application"`
				}) (commsImpl.IConnectionReactorFactory, error) {
					return listener.NewConnectionReactorFactory(
						TextListenerConnection,
						params.PubSub,
						func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
							bytes, err := json.MarshalIndent(m, "", "\t")
							if err != nil {
								return nil, err
							}
							return multiBlock.NewReaderWriterBlock(bytes), nil
						}), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: commsImpl.NewNetListenApp(
					TextListenerConnection,
					url,
					commsImpl.TransportFactoryEmptyName,
					TextListenerConnection,
					pairInformation),
			}),
	)
}

func CompressedListener(url string, pairInformation ...*common.PairInformation) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: commsImpl.ConnectionReactorFactoryConst,
				Target: func(params struct {
					fx.In
					PubSub *pubsub.PubSub `name:"Application"`
				}) (commsImpl.IConnectionReactorFactory, error) {
					return listener.NewConnectionReactorFactory(
						CompressedListenerConnection,
						params.PubSub,
						func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
							return stream.Marshall(data)
						}), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: commsImpl.NewNetListenApp(
					CompressedListenerConnection,
					url,
					commsImpl.TransportFactoryCompressedName,
					CompressedListenerConnection,
					pairInformation),
			}),
	)
}
