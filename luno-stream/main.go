package main

import (
	"fmt"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/lunoStream"
	"github.com/bhbosman/gocommon/model"
)

func main() {
	var errors []error
	dd := func(pairName string, useProxy bool) (*common.PairInformation, error) {
		return common.NewPairInformation(
			pairName,
			common.LunoDialerXBTZARServiceNumber,
			model.NoDependency,
			useProxy,
			"tcp4://127.0.0.1:1080",
			fmt.Sprintf("wss://ws.luno.com:443/api/1/stream/%v", pairName),
		)
	}
	useProxy := false
	pairInfoXBTZAR, err := dd("XBTZAR", useProxy)
	if err != nil {
		errors = append(errors, err)
	}
	//pairInfoXBTEUR, err := dd("XBTEUR", useProxy)
	//if err != nil {
	//	errors = append(errors, err)
	//}
	//pairInfoXBTUGX, err := dd("XBTUGX", useProxy)
	//if err != nil {
	//	errors = append(errors, err)
	//}
	//pairInfoXBTZMW, err := dd("XBTZMW", useProxy)
	//if err != nil {
	//	errors = append(errors, err)
	//}
	//pairInfoETHXBT, err := dd("ETHXBT", useProxy)
	//if err != nil {
	//	errors = append(errors, err)
	//}
	//pairInfoBCHXBT, err := dd("BCHXBT", useProxy)
	//if err != nil {
	//	errors = append(errors, err)
	//}

	lunoStreamApp, err := lunoStream.App(
		false,
		lunoStream.NewErrors(errors...),
		lunoStream.TextListenerUrl("tcp4://127.0.0.1:3000", common.TextListenerServiceNumber, model.NoDependency),
		lunoStream.CompressedListenerUrl("tcp4://127.0.0.1:3001", common.CompressedListenerServiceNumber, model.NoDependency),
		lunoStream.AddCurrencyPair(pairInfoXBTZAR),
		//lunoStream.AddCurrencyPair(pairInfoETHXBT),
		//lunoStream.AddCurrencyPair(pairInfoXBTUGX),
		//lunoStream.AddCurrencyPair(pairInfoXBTEUR),
		//lunoStream.AddCurrencyPair(pairInfoXBTZMW),
		//lunoStream.AddCurrencyPair(pairInfoBCHXBT),
	)
	if err != nil {

	}
	if lunoStreamApp.FxApp.Err() != nil {
		println(lunoStreamApp.FxApp.Err().Error())
		return
	}
	lunoStreamApp.RunTerminalApp()
}
