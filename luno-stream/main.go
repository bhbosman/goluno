package main

import (
	"flag"
	"fmt"
	"github.com/bhbosman/goLuno/internal/flags"
	"github.com/bhbosman/goLuno/internal/lunoStream"
)

func main() {
	flag.Parse()
	lunoStreamApp, err := lunoStream.App(
		false,
		*flags.ApplicationName,
		lunoStream.CompressedListenerUrl(fmt.Sprintf("tcp4://127.0.0.1:3001")),
	)
	if err != nil {

	}
	if lunoStreamApp.FxApp.Err() != nil {
		println(lunoStreamApp.FxApp.Err().Error())
		return
	}
	lunoStreamApp.RunTerminalApp()
}
