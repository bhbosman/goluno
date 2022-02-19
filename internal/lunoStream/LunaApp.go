package lunoStream

import (
	"github.com/rivo/tview"
	"go.uber.org/fx"
)

type LunaApp struct {
	FxApp         *fx.App
	ShutDowner    fx.Shutdowner
	UiApplication *tview.Application
}
