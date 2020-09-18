package main

import (
	"github.com/bhbosman/goLuno/internal"
	"github.com/bhbosman/goLuno/internal/common"
	"time"
)

func main() {
	app, _ := internal.LunoStreamApp(
		internal.HttpListenerUrl("http://127.0.0.1:8080"),
		internal.TextListenerUrl("tcp4://127.0.0.1:3000"),
		internal.CompressedListenerUrl("tcp4://127.0.0.1:3001"),
		internal.AddCurrencyPair(common.NewPairInformation("XBTZAR")),
		internal.AddCurrencyPair(common.NewPairInformation("XBTEUR")),
		internal.AddCurrencyPair(common.NewPairInformation("XBTUGX")),
		internal.AddCurrencyPair(common.NewPairInformation("XBTZMW")),
		internal.AddCurrencyPair(common.NewPairInformation("ETHXBT")),
		internal.AddCurrencyPair(common.NewPairInformation("BCHXBT")))
	if app.Err() != nil {
		return
	}
	app.Run()
	// allow shutdown to complete
	time.Sleep(time.Second)
}
