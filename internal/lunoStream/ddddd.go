package lunoStream

import (
	"github.com/bhbosman/gocommon/ui/uiImpl"
	"github.com/bhbosman/gocommon/ui/uiIntf"
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

type ICommand interface {
	Callback(app *tview.Application) func()
	MainText() string
	SecondaryText() string
	ShortCut() rune
}
type Command struct {
	mainText      string
	secondaryText string
	shortCut      rune
	cb            func(app *tview.Application) func()
}

func NewCommand(mainText string, secondaryText string, shortCut rune, cb func(app *tview.Application) func()) *Command {
	return &Command{mainText: mainText, secondaryText: secondaryText, shortCut: shortCut, cb: cb}
}

func (self *Command) ShortCut() rune {
	return self.shortCut
}

func (self *Command) SecondaryText() string {
	return self.secondaryText
}

func (self *Command) MainText() string {
	return self.mainText
}

func (self *Command) Callback(app *tview.Application) func() {
	return self.cb(app)
}

func terminalApplicationOptionsss() []fx.Option {
	return []fx.Option{
		fx.Provide(
			fx.Annotated{
				Target: func() *TerminalApplicationBuilder {
					return &TerminalApplicationBuilder{}
				},
			}),
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
					TerminalApplicationBuilder *TerminalApplicationBuilder
					MainPages                  *tview.Pages `name:"MainPages"`
					MainPageCommandList        []ICommand   `group:"MainPageCommandList"`
					UiApp                      uiIntf.IUiService
				}) *tview.Application {
					return params.UiApp.Build()()

				},
			}),
	}
}
