package internal

import (
	lunoInternal "github.com/bhbosman/goLuno/internal"
	"github.com/bhbosman/goLuno/internal/common"
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
	self.app, self.shutDowner = lunoInternal.LunoStreamApp(
		lunoInternal.HttpListenerUrl("http://127.0.0.1:8080"),
		lunoInternal.TextListenerUrl("tcp4://127.0.0.1:3000"),
		lunoInternal.CompressedListenerUrl("tcp4://127.0.0.1:3001"),
		lunoInternal.AddCurrencyPair(common.NewPairInformation("XBTZAR")),
		lunoInternal.AddCurrencyPair(common.NewPairInformation("XBTEUR")),
		lunoInternal.AddCurrencyPair(common.NewPairInformation("XBTUGX")),
		lunoInternal.AddCurrencyPair(common.NewPairInformation("XBTZMW")),
		lunoInternal.AddCurrencyPair(common.NewPairInformation("ETHXBT")),
		lunoInternal.AddCurrencyPair(common.NewPairInformation("BCHXBT")))
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
