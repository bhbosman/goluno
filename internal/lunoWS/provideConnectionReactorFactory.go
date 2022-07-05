package lunoWS

import (
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func ProvideConnectionReactorFactory() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					LunoAPIKeyID     string         `name:"LunoAPIKeyID"`
					LunoAPIKeySecret string         `name:"LunoAPIKeySecret"`
					PubSub           *pubsub.PubSub `name:"Application"`
				},
			) (intf.IConnectionReactorFactory, error) {
				fac := NewConnectionReactorFactory(
					"",
					params.LunoAPIKeyID,
					params.LunoAPIKeySecret,
					params.PubSub)
				return fac, nil
			},
		},
	)
}
