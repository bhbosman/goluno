package listener

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetListener"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/bvisMessageBreaker"
	"github.com/bhbosman/goCommsStacks/messageCompressor"
	"github.com/bhbosman/goCommsStacks/messageNumber"
	"github.com/bhbosman/goCommsStacks/pingPong"
	"github.com/bhbosman/goCommsStacks/protoBuf"
	"github.com/bhbosman/goCommsStacks/topStack"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"net/url"
)

func CompressedListener(
	maxConnections int,
	UseProxy bool,
	ProxyUrl *url.URL,
	ConnectionUrl *url.URL,
) fx.Option {
	const CompressedListenerConnection = "Luno Client Connection Manager(TlS Compressed)"

	return fx.Provide(
		fx.Annotated{
			Group: "Apps",
			Target: func(
				params struct {
					fx.In
					PubSub               *pubsub.PubSub `name:"Application"`
					NetAppFuncInParams   common2.NetAppFuncInParams
					FullMarketDataHelper fullMarketDataHelper.IFullMarketDataHelper
					FmdService           fullMarketDataManagerService.IFmdManagerService
				},
			) messages.CreateAppCallback {
				f := goCommsNetListener.NewNetListenApp(
					CompressedListenerConnection,
					CompressedListenerConnection,
					UseProxy,
					ProxyUrl,
					ConnectionUrl,
					common2.MaxConnectionsSetting(maxConnections),
					common2.NewConnectionInstanceOptions(
						fx.Provide(
							fx.Annotated{
								Target: func() (fullMarketDataHelper.IFullMarketDataHelper, fullMarketDataManagerService.IFmdManagerService) {
									return params.FullMarketDataHelper, params.FmdService
								},
							},
						),
						PubSub.ProvidePubSubInstance("Application", params.PubSub),
						goCommsDefinitions.ProvideTransportFactoryForCompressedName(
							topStack.Provide(),
							pingPong.Provide(),
							protoBuf.Provide(),
							messageCompressor.Provide(),
							messageNumber.Provide(),
							bvisMessageBreaker.Provide(),
							bottom.Provide(),
						),
						ProvideConnectionReactor(),
					),
				)
				return f(
					params.NetAppFuncInParams)
			},
		},
	)
}
