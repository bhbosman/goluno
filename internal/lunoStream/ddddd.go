package lunoStream

import (
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"golang.org/x/net/context"
	"time"
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
		fx.Provide(
			fx.Annotated{
				Group: "MainPageCommandList",
				Target: func(params struct {
					fx.In
					Connections connectionManager.IConnectionManagerService
					MainPages   *tview.Pages `name:"MainPages"`
				}) ICommand {
					return NewCommand("Connections", "", 0,
						func(app *tview.Application) func() {
							return func() {
								connectionList, _ := params.Connections.GetConnections(context.Background())
								commandList := tview.NewList()
								commandList.
									ShowSecondaryText(false).
									SetBorder(true).
									SetTitle("Commands")
								commandList.AddItem("..", "", 0, func() {
									params.MainPages.RemovePage("Connections")
								})

								commandList.AddItem("All", "", 0, func() {
									commandList.ShowSecondaryText(false)
								})
								for _, enumValue := range connectionList {
									information := enumValue
									//index := i
									commandList.AddItem(information.Name, "", 0, func() {

									})
								}
								table := tview.NewTable()
								dd := tview.NewFlex().
									AddItem(
										tview.NewFlex().
											SetDirection(tview.FlexRow).
											AddItem(commandList, 10, 1, true).
											AddItem(table, 0, 1, false),
										0,
										1,
										true)
								params.MainPages.AddPage("Connections", dd, true, true)
							}
						})
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "MainPageCommandList",
				Target: func() ICommand {
					return NewCommand("Services", "", 0,
						func(app *tview.Application) func() {
							return func() {

							}
						})
				},
			}),

		fx.Provide(
			fx.Annotated{
				Target: func(params struct {
					fx.In
					TerminalApplicationBuilder *TerminalApplicationBuilder
					MainPages                  *tview.Pages `name:"MainPages"`
					MainPageCommandList        []ICommand   `group:"MainPageCommandList"`
				}) *tview.Application {

					app := tview.NewApplication()

					commandList := params.TerminalApplicationBuilder.createCommandList()
					for _, command := range params.MainPageCommandList {
						commandList.AddItem(command.MainText(), command.SecondaryText(), command.ShortCut(), command.Callback(app))
					}
					commandList.AddItem("Quit", "", 'q', func() {
						app.Stop()
					})

					outputPanel := params.TerminalApplicationBuilder.createOutputPanel(app)

					timeText := tview.NewTextView().SetTextAlign(tview.AlignRight)
					timeText.SetText(time.Now().Format(time.Stamp))
					layout := params.TerminalApplicationBuilder.createMainLayout(commandList, outputPanel, timeText)
					params.MainPages.AddPage("main", layout, true, true)

					app.SetRoot(params.MainPages, true)

					return app
				},
			}),
	}
}
