package internal

import (
	"github.com/bhbosman/goFxApp"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/lunoStream"
	"github.com/kardianos/service"
)

type Program struct {
	app *goFxApp.TerminalAppUsingFxApp
}

func NewProgram() *Program {
	return &Program{
		app: nil,
	}
}

func (self *Program) Start(_ service.Service) error {
	self.app, _ = lunoStream.App(
		true,
		lunoStream.HttpListenerUrl("http://127.0.0.1:8080"),
		lunoStream.TextListenerUrl("tcp4://127.0.0.1:3000"),
		lunoStream.CompressedListenerUrl("tcp4://127.0.0.1:3001"),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTZAR")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTEUR")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTUGX")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTZMW")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("ETHXBT")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("BCHXBT")))
	if self.app.FxApp.Err() != nil {
		return self.app.FxApp.Err()
	}

	go func() {
		self.app.FxApp.Run()
	}()
	return nil
}

func (self *Program) Stop(_ service.Service) error {
	return self.app.Shutdown.Shutdown()
}
