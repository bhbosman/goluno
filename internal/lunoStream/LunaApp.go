package lunoStream

import "go.uber.org/fx"

type LunaApp struct {
	FxApp      *fx.App
	ShutDowner fx.Shutdowner
}
