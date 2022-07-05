package lunoStream

import (
	"encoding/json"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goCommsNetListener"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/top"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/listener"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"net/url"
)

func TextListener(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl,
	maxConnections int,
	UseProxy bool,
	ProxyUrl *url.URL,
	ConnectionUrl *url.URL,
	pairInformation ...*common.PairInformation) fx.Option {
	const TextListenerConnection = "TextListenerConnection"

	return fx.Provide(
		fx.Annotated{
			Group: "Apps",
			Target: func(params struct {
				fx.In
				PubSub             *pubsub.PubSub `name:"Application"`
				NetAppFuncInParams common2.NetAppFuncInParams
			}) (messages.CreateAppCallback, error) {

				f := goCommsNetListener.NewNetListenApp(
					TextListenerConnection,
					serviceIdentifier,
					serviceDependentOn,
					TextListenerConnection,
					UseProxy,
					ProxyUrl,
					ConnectionUrl,
					goCommsDefinitions.TransportFactoryEmptyName,
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
						goCommsDefinitions.ProvideTransportFactoryForEmptyName(
							top.ProvideTopStack(),
							bottom.Provide(),
						),
						ProvideConnectionReactorFactory(),
					),
				)
				return f(params.NetAppFuncInParams), nil
			},
		})
}

func ProvideConnectionReactorFactory() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					PubSub          *pubsub.PubSub `name:"Application"`
					ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl
				}) (intf.IConnectionReactorFactory, error) {
				return listener.NewConnectionReactorFactory(
					params.PubSub,
					func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
						bytes, err := json.MarshalIndent(m, "", "\t")
						if err != nil {
							return nil, err
						}
						return gomessageblock.NewReaderWriterBlock(bytes), nil
					},
					params.ConsumerCounter)
			},
		},
	)
}
