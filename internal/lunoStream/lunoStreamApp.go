package lunoStream

import (
	"github.com/bhbosman/goCommonMarketData/fullMarketDataHelper"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerService"
	"github.com/bhbosman/goCommonMarketData/fullMarketDataManagerViewer"
	"github.com/bhbosman/goCommonMarketData/instrumentReference"
	"github.com/bhbosman/goCommsMultiDialer"
	"github.com/bhbosman/goFxApp"
	"github.com/bhbosman/goLuno/internal/lunoConfiguration"
	app2 "github.com/bhbosman/gocommon/Providers"
	"go.uber.org/fx"
	"go.uber.org/multierr"
)

func App(
	serviceApplication bool,
	applicationName string,
	pairs ...ILunoStreamAppApplySettings,
) (*goFxApp.TerminalAppUsingFxApp, error) {

	settings := &AppSettings{
		textListenerEnabled:       false,
		textListenerUrl:           "tcp4://127.0.0.1:3000",
		compressedListenerEnabled: false,
		compressedListenerUrl:     "tcp4://127.0.0.1:3001",
		canDial:                   nil,
		macConnections:            0,
		errors:                    nil,
	}
	var errs error
	for _, apply := range pairs {
		err := apply.apply(settings)
		errs = multierr.Append(errs, err)
	}
	if errs != nil {
		return nil, errs
	}
	var runTimeManager *app2.RunTimeManager

	return goFxApp.NewFxMainApplicationServices(
			applicationName,
			serviceApplication,
			fx.Error(settings.errors...),
			fx.Supply(settings),
			fx.Populate(&runTimeManager),
			goCommsMultiDialer.Provide(),
			fullMarketDataManagerViewer.Provide(),
			fullMarketDataManagerService.Provide(false),
			fullMarketDataHelper.Provide(),
			instrumentReference.Provide(),
			app2.RegisterRunTimeManager(),
			lunoConfiguration.Provide(),
			ProvideCompressedListener(settings),
			ProvideLunoKeys(),
			//true,
			//&lunoKeys{
			//	Key:    "e52n78axhy2j7",
			//	Secret: "4q00paAkXche01noiISYWsZQGtSOKe1kuMnQUk3m3Io",
			//}),
			ProvideLunoAPIKeyID(),
			ProvideLunoAPIKeySecret(),
		),
		nil
}
