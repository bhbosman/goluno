package lunoWS

import (
	"fmt"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func ProvideDialers(options ...DialersApply) fx.Option {
	settings := &lunoStreamDialersSettings{}
	for _, option := range options {
		if option == nil {
			continue
		}
		option.apply(settings)
	}
	var opt []fx.Option
	for _, option := range settings.pairs {
		if option == nil {
			continue
		}
		crfName := fmt.Sprintf("LunoStream.%v.CRF", option.Pair)
		opt = append(
			opt,
			fx.Provide(
				fx.Annotated{
					Group: "CFR",
					Target: func(params struct {
						fx.In
						PubSub           *pubsub.PubSub `name:"Application"`
						LunoAPIKeyID     string         `name:"LunoAPIKeyID"`
						LunoAPIKeySecret string         `name:"LunoAPIKeySecret"`
					}) (intf.IConnectionReactorFactory, error) {
						cfr := NewConnectionReactorFactory(
							crfName,
							params.LunoAPIKeyID,
							params.LunoAPIKeySecret,
							params.PubSub)
						return cfr, nil
					},
				},
			))
		opt = append(
			opt,
			fx.Provide(
				fx.Annotated{
					Group: "Apps",
					Target: netDial.NewNetDialApp(
						fmt.Sprintf("luno stream[%v]", option.Pair),
						fmt.Sprintf("wss://ws.luno.com:443/api/1/stream/%v", option.Pair),
						impl.WebSocketName,
						crfName,
						netDial.MaxConnectionsSetting(settings.maxConnections),
						netDial.UserContextValue(option),
						netDial.CanDial(settings.canDial...)),
				}))
	}
	return fx.Options(opt...)
}
