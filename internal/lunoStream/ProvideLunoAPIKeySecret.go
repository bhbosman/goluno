package lunoStream

import "go.uber.org/fx"

func ProvideLunoAPIKeySecret() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "LunoAPIKeySecret",
			Target: func(data *lunoKeys) string {
				return data.Secret
			},
		})
}
