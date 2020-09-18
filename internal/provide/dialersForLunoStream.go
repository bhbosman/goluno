package provide

import (
	"fmt"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/lunoWS"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func LunoStreamDialers(lunoPairInformation ...*common.PairInformation) fx.Option {
	const LunoStreamConnectionReactorFactory = "LunoStreamConnectionReactorFactory"
	var opt []fx.Option
	opt = append(opt, fx.Provide(fx.Annotated{
		Group: commsImpl.ConnectionReactorFactoryConst,
		Target: func(
			params struct {
				fx.In
				PubSub       *pubsub.PubSub `name:"Application"`
				APIKeyID     string         `name:"LunoAPIKeyID"`
				APIKeySecret string         `name:"LunoAPIKeySecret"`
			}) (commsImpl.IConnectionReactorFactory, error) {

			return lunoWS.NewConnectionReactorFactory(
				LunoStreamConnectionReactorFactory,
				params.APIKeyID,
				params.APIKeySecret,
				params.PubSub), nil

		},
	}))
	for _, information := range lunoPairInformation {
		opt = append(opt, fx.Provide(fx.Annotated{
			Group: "Apps",
			Target: commsImpl.NewNetDialApp(
				fmt.Sprintf("luno stream[%v]", information.Pair),
				fmt.Sprintf("wss://ws.luno.com:443/api/1/stream/%v", information.Pair),
				commsImpl.WebSocketName,
				LunoStreamConnectionReactorFactory,
				information),
		}))
	}
	return fx.Options(opt...)
}
