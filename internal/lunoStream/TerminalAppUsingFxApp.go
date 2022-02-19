package lunoStream

import (
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"golang.org/x/net/context"
	"os"
)

type TerminalAppUsingFxApp struct {
	FxApp       *fx.App
	TerminalApp *tview.Application
	Shutdown    fx.Shutdowner
}

func NewTerminalAppUsingFxApp(fxApp *fx.App, terminalApp *tview.Application, shutdown fx.Shutdowner) *TerminalAppUsingFxApp {
	return &TerminalAppUsingFxApp{
		FxApp:       fxApp,
		TerminalApp: terminalApp,
		Shutdown:    shutdown,
	}
}

func (self *TerminalAppUsingFxApp) RunTerminalApp() {
	startCtx, cancel := context.WithTimeout(context.Background(), self.FxApp.StartTimeout())
	err := self.FxApp.Start(startCtx)
	if err != nil {
		os.Exit(1)
	}
	cancel()
	if err := self.TerminalApp.Run(); err != nil {
		panic(err)
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), self.FxApp.StopTimeout())

	err = self.FxApp.Stop(stopCtx)
	defer cancel()

}
