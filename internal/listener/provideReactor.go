package listener

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goprotoextra"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
)

func ProvideConnectionReactor() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						Logger                 *zap.Logger
						CancelCtx              context.Context
						CancelFunc             context.CancelFunc
						ConnectionCancelFunc   model.ConnectionCancelFunc
						PubSub                 *pubsub.PubSub `name:"Application"`
						GoFunctionCounter      GoFunctionCounter.IService
						UniqueReferenceService interfaces.IUniqueReferenceService
						FullMarketDataHelper   fullMarketDataHelper.IFullMarketDataHelper
						FmdService             fullMarketDataManagerService.IFmdManagerService
					},
				) (intf.IConnectionReactor, error) {
					params.Logger.Info("Creating Connection Reactor")
					result, err := NewConnectionReactor(
						params.Logger,
						params.CancelCtx,
						params.CancelFunc,
						params.ConnectionCancelFunc,
						params.PubSub,
						func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
							return stream.Marshall(data)
						},
						params.GoFunctionCounter,
						params.UniqueReferenceService,
						params.FullMarketDataHelper,
						params.FmdService,
					)
					if err != nil {
						return nil, err
					}
					return result, nil
				},
			},
		),
	)
}
