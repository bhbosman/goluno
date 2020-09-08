package main

import (
	"context"
	"github.com/bhbosman/goLunoApi/register"
	"github.com/bhbosman/goLunoApi/stream"
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/stacks/connectionManager"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.LogName("Luno Stream Application"),
		app2.RegisterRootContext(),
		connectionManager.RegisterDefaultConnectionManager(),
		commsImpl.RegisterAllConnectionRelatedServices(),
		register.DialerForLunoStream(
			stream.NewLunoPairInformation("XBTZAR"),
			stream.NewLunoPairInformation("XBTEUR"),
			stream.NewLunoPairInformation("XBTUGX"),
			stream.NewLunoPairInformation("XBTZMW"),
			stream.NewLunoPairInformation("ETHXBT"),
			stream.NewLunoPairInformation("BCHXBT")),
		register.ReadLunoKeys(),
		fx.Invoke(
			func(params struct {
				fx.In
				Lifecycle      fx.Lifecycle
				Apps           []*fx.App `group:"Apps"`
				Logger         fx.ILogger
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
			}),
	)
	if app.Err() != nil {
		return
	}
	app.Run()
}
