package main

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/lunoStream"
)

func main() {
	lunoStreamApp, err := lunoStream.App(false,
		lunoStream.HttpListenerUrl("http://127.0.0.1:8080"),
		lunoStream.TextListenerUrl("tcp4://127.0.0.1:3000"),
		lunoStream.CompressedListenerUrl("tcp4://127.0.0.1:3001"),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTZAR")),
		//lunoStream.AddCurrencyPair(common.NewPairInformation("XBTEUR")),
		//lunoStream.AddCurrencyPair(common.NewPairInformation("XBTUGX")),
		//lunoStream.AddCurrencyPair(common.NewPairInformation("XBTZMW")),
		//lunoStream.AddCurrencyPair(common.NewPairInformation("ETHXBT")),
		//lunoStream.AddCurrencyPair(common.NewPairInformation("BCHXBT")),
	)
	if err != nil {

	}
	if lunoStreamApp.FxApp.Err() != nil {
		println(lunoStreamApp.FxApp.Err().Error())
		return
	}
	lunoStreamApp.RunTerminalApp()
}
