package reactorWS

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func Provide() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					Pair                   string         `name:"Pair"`
					LunoAPIKeyID           string         `name:"LunoAPIKeyID"`
					LunoAPIKeySecret       string         `name:"LunoAPIKeySecret"`
					PubSub                 *pubsub.PubSub `name:"Application"`
					CancelCtx              context.Context
					CancelFunc             context.CancelFunc
					ConnectionCancelFunc   model.ConnectionCancelFunc
					Logger                 *zap.Logger
					UniqueReferenceService interfaces.IUniqueReferenceService
					GoFunctionCounter      GoFunctionCounter.IService
					FmdService             fullMarketDataManagerService.IFmdManagerService
					FmdHelper              fullMarketDataHelper.IFullMarketDataHelper
					ReferenceData          instrumentReference.LunoReferenceData
				},
			) (intf.IConnectionReactor, error) {
				result, err := NewConnectionReactor(
					params.Logger,
					params.CancelCtx,
					params.CancelFunc,
					params.ConnectionCancelFunc,
					params.LunoAPIKeyID,
					params.LunoAPIKeySecret,
					params.PubSub,
					params.ReferenceData,
					params.GoFunctionCounter,
					params.UniqueReferenceService,
					params.FmdHelper,
					params.FmdService,
				)
				if err != nil {
					return nil, err
				}
				return result, nil
			},
		},
	)
}
