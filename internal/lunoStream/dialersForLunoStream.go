package lunoStream

import (
	"fmt"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/lunoWS"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/netDial"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func Dialers(options ...DialersApply) fx.Option {
	settings := &lunoStreamDialersSettings{}
	for _, option := range options {
		option.apply(settings)
	}

	const LunoStreamConnectionReactorFactory = "LunoStreamConnectionReactorFactory"
	var opt []fx.Option
	opt = append(opt, fx.Provide(fx.Annotated{
		Group: commsImpl.ConnectionReactorFactoryConst,
		Target: func(
			params struct {
				fx.In
				PubSub *pubsub.PubSub `name:"Application"`
				APIKeyID        string `name:"LunoAPIKeyID"`
				APIKeySecret    string `name:"LunoAPIKeySecret"`
			}) (commsImpl.IConnectionReactorFactory, error) {

			return lunoWS.NewConnectionReactorFactory(
				LunoStreamConnectionReactorFactory,
				params.APIKeyID,
				params.APIKeySecret,
				params.PubSub), nil

		},
	}))
	for _, option := range settings.pairs {
		opt = append(opt, fx.Provide(fx.Annotated{
			Group: "Apps",
			Target: netDial.NewNetDialApp(
				fmt.Sprintf("luno stream[%v]", option.Pair),
				fmt.Sprintf("wss://ws.luno.com:443/api/1/stream/%v", option.Pair),
				commsImpl.WebSocketName,
				LunoStreamConnectionReactorFactory,
				netDial.UserContextValue(option),
				netDial.CanDial(settings.canDial...)),
		}))
	}
	return fx.Options(opt...)
}

type lunoStreamDialersSettings struct {
	pairs   []*common.PairInformation
	canDial []netDial.ICanDial
}
type DialersApply interface {
	apply(*lunoStreamDialersSettings)
}
type addPairsInformation struct {
	pairs []*common.PairInformation
}

func AddPairsInformation(pairs []*common.PairInformation) *addPairsInformation {
	return &addPairsInformation{pairs: pairs}
}

func (self addPairsInformation) apply(settings *lunoStreamDialersSettings) {
	for _, pair := range self.pairs {
		settings.pairs = append(settings.pairs, pair)
	}
}

type canDialSetting struct {
	canDial []netDial.ICanDial
}

func CanDial(canDial ...netDial.ICanDial) *canDialSetting {
	return &canDialSetting{canDial: canDial}
}

func (self canDialSetting) apply(settings *lunoStreamDialersSettings) {
	for _, cd := range self.canDial {
		settings.canDial = append(settings.canDial, cd)
	}
}
