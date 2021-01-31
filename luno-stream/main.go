package main

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/lunoStream"
	"time"
)

func main() {
	app, _ := lunoStream.App(
		lunoStream.HttpListenerUrl("http://127.0.0.1:8080"),
		lunoStream.TextListenerUrl("tcp4://127.0.0.1:3000"),
		lunoStream.CompressedListenerUrl("tcp4://127.0.0.1:3001"),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTZAR")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTEUR")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTUGX")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("XBTZMW")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("ETHXBT")),
		lunoStream.AddCurrencyPair(common.NewPairInformation("BCHXBT")))
	if app.Err() != nil {
		println(app.Err().Error())
		return
	}
	app.Run()
	// allow shutdown to complete
	time.Sleep(time.Second)
}
