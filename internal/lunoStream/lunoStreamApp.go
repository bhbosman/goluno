package lunoStream

import (
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goFxApp"
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/lunoWS"
	app2 "github.com/bhbosman/gocommon/Providers"
	"github.com/bhbosman/gocommon/model"
	"go.uber.org/fx"
	"go.uber.org/multierr"
)

func App(serviceApplication bool, pairs ...ILunoStreamAppApplySettings) (*goFxApp.TerminalAppUsingFxApp, error) {
	settings := &AppSettings{
		pairs:                     nil,
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
	ConsumerCounter := goCommsNetDialer.NewCanDialDefaultImpl()
	var runTimeManager *app2.RunTimeManager

	return goFxApp.NewFxMainApplicationServices(
		"LunoApplication",
		serviceApplication,
		fx.Error(settings.errors...),
		fx.Supply(settings, ConsumerCounter),
		fx.Populate(&runTimeManager),
		app2.RegisterRunTimeManager(),
		ProvideTextListener(common.TextListenerServiceNumber, model.NoDependency, settings, ConsumerCounter),
		ProvideCompressedListener(common.CompressedListenerServiceNumber, model.NoDependency, settings, ConsumerCounter),
		lunoWS.ProvideDialers(
			lunoWS.CanDial(ConsumerCounter),
			lunoWS.AddPairsInformation(settings.pairs),
			lunoWS.MaxConnections(1)),
		ProvideLunoKeys(
			false,
			&lunoKeys{
				Key:    "hzy4572ygxbb6",
				Secret: "0LXalWARHJmhze3Yk0lUPrHF51lUn8XqX49E4D7vsW4",
			}),
		ProvideLunoAPIKeyID(),
		ProvideLunoAPIKeySecret(),
	), nil
}
