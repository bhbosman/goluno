package internal

import (
	"context"
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
)

func InvokeApps() fx.Option {
	return fx.Invoke(
		func(params struct {
			fx.In
			Lifecycle      fx.Lifecycle
			Apps           []*fx.App `group:"Apps"`
			Logger         *gologging.Factory
			RunTimeManager *app2.RunTimeManager
		}) {
			params.Lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return params.RunTimeManager.Start(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return params.RunTimeManager.Stop(ctx)
				},
			})
			for _, item := range params.Apps {
				localApp := item
				params.Lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return localApp.Start(ctx)
					},
					OnStop: func(ctx context.Context) error {
						return localApp.Stop(ctx)
					},
				})
			}
		})
}
