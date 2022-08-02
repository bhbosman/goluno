package lunoStream

import "go.uber.org/fx"

func ProvideLunoAPIKeyID() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "LunoAPIKeyID",
			Target: func(data *lunoKeys) string {
				return data.Key
			},
		})
}
