package lunoStream

import (
	"github.com/bhbosman/gocommon/ui/uiImpl"
	"github.com/bhbosman/gocommon/ui/uiIntf"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

func terminalApplicationOptionsss() []fx.Option {
	return []fx.Option{
		fx.Provide(
			fx.Annotated{
				Name: "MainPages",
				Target: func() *tview.Pages {
					return tview.NewPages()
				}}),

		fx.Provide(fx.Annotated{
			Target: func(params struct {
				fx.In
				ApplicationContext context.Context `name:"Application"`
				PubSub             *pubsub.PubSub  `name:"Application"`
			}) uiIntf.IUiService {
				return uiImpl.NewService(params.ApplicationContext, params.PubSub, tview.NewApplication())
			},
		}),
		fx.Provide(
			fx.Annotated{
				Target: func(params struct {
					fx.In
					UiApp uiIntf.IUiService
				}) *tview.Application {
					return params.UiApp.Build()()
				},
			}),
	}
}
