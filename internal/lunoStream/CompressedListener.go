package lunoStream

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsNetListener"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/bvisMessageBreaker"
	"github.com/bhbosman/goCommsStacks/messageCompressor"
	"github.com/bhbosman/goCommsStacks/messageNumber"
	"github.com/bhbosman/goCommsStacks/pingPong"
	"github.com/bhbosman/goCommsStacks/tlsConnection"
	"github.com/bhbosman/goCommsStacks/top"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"net/url"
)

func CompressedListener(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl,
	maxConnections int,
	UseProxy bool,
	ProxyUrl *url.URL,
	ConnectionUrl *url.URL,
	pairInformation ...*common.PairInformation,
) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"

	return fx.Provide(
		fx.Annotated{
			Group: "Apps",
			Target: func(params struct {
				fx.In
				PubSub             *pubsub.PubSub `name:"Application"`
				NetAppFuncInParams common2.NetAppFuncInParams
			}) messages.CreateAppCallback {
				f := goCommsNetListener.NewNetListenApp(
					CompressedListenerConnection,
					serviceIdentifier,
					serviceDependentOn,
					CompressedListenerConnection,
					UseProxy,
					ProxyUrl,
					ConnectionUrl,
					goCommsDefinitions.TransportFactoryCompressedTlsName,
					goCommsNetListener.UserContextValue(pairInformation),
					common2.MaxConnectionsSetting(maxConnections),
					common2.NewConnectionInstanceOptions(
						fx.Provide(
							fx.Annotated{
								Target: func() *goCommsNetDialer.CanDialDefaultImpl {
									return ConsumerCounter
								},
							},
						),
						PubSub.ProvidePubSubInstance("Application", params.PubSub),
						goCommsDefinitions.ProvideTransportFactoryForCompressedTlsName(
							top.ProvideTopStack(),
							pingPong.ProvidePingPongStacks(),
							messageCompressor.Provide(),
							messageNumber.ProvideMessageNumberStack(),
							bvisMessageBreaker.Provide(),
							tlsConnection.ProvideTlsConnectionStacks(),
							bottom.Provide(),
						),
						ProvideConnectionReactorFactory2(),
					),
				)
				return f(
					params.NetAppFuncInParams)
			},
		},
	)
}

func ProvideConnectionReactorFactory2() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					PubSub          *pubsub.PubSub `name:"Application"`
					ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl
				},
			) (intf.IConnectionReactorFactory, error) {
				return listener.NewConnectionReactorFactory(
					params.PubSub,
					func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
						return stream.Marshall(data)
					},
					params.ConsumerCounter)
			},
		},
	)
}
