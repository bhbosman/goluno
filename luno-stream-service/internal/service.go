package internal

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/lunoStream"
	"github.com/kardianos/service"
	"go.uber.org/fx"
)

type Program struct {
	app        *fx.App
	shutDowner fx.Shutdowner
}

func NewProgram() *Program {
	return &Program{
		app:        nil,
		shutDowner: nil,
	}
}

func (self *Program) Start(s service.Service) error {
	self.app, self.shutDowner = lunoStream.App(
		lunoStream.HttpListenerUrl("http://127.0.0.1:8080"),
		lunoStream.TextListenerUrl("tcp4://127.0.0.1:3000"),
		lunoStream.CompressedListenerUrl("tcp4://127.0.0.1:3001"),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTZAR")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTEUR")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTUGX")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTZMW")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("ETHXBT")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("BCHXBT")))
	if self.app.Err() != nil {
		return self.app.Err()
	}

	go func() {
		self.app.Run()
	}()
	return nil

}

func (self *Program) Stop(s service.Service) error {
	return self.shutDowner.Shutdown()
}
