package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/lunoWS"
	"github.com/bhbosman/gocommon"
	app2 "github.com/bhbosman/gocommon/Providers"
	"github.com/bhbosman/gocommon/ui/uiImpl"
	"github.com/bhbosman/gocomms/connectionManager/CMImpl"
	"github.com/bhbosman/gocomms/connectionManager/endpoints"
	"github.com/bhbosman/gocomms/connectionManager/view"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/provide"
	"github.com/rivo/tview"
	"go.uber.org/fx"
	"go.uber.org/multierr"
)

func App(serviceApplication bool, pairs ...ILunoStreamAppApplySettings) (*TerminalAppUsingFxApp, error) {
	settings := &AppSettings{
		pairs:                     nil,
		textListenerEnabled:       false,
		textListenerUrl:           "tcp4://127.0.0.1:3000",
		compressedListenerEnabled: false,
		compressedListenerUrl:     "tcp4://127.0.0.1:3001",
		httpListenerUrlEnabled:    false,
		httpListenerUrl:           "http://127.0.0.1:8080",
		canDial:                   nil,
		macConnections:            0,
	}
	var errs error
	for _, apply := range pairs {
		err := apply.apply(settings)
		errs = multierr.Append(errs, err)
	}
	if errs != nil {
		return nil, errs
	}
	ConsumerCounter := netDial.NewCanDialDefaultImpl()
	var shutDowner fx.Shutdowner
	var runTimeManager *gocommon.RunTimeManager
	var terminalApplication *tview.Application

	terminalApplicationOptions := fx.Options()
	if !serviceApplication {
		ssss := uiImpl.TerminalApplicationOptionsss()
		options := append([]fx.Option{fx.Populate(&terminalApplication)}, ssss...)
		terminalApplicationOptions = fx.Options(options...)
	}

	fxApp := app2.NewFxAppWithServices(
		"LunoApplication",
		false,
		terminalApplicationOptions,
		fx.Supply(settings, ConsumerCounter),
		fx.Populate(&shutDowner),
		fx.Populate(&runTimeManager),
		app2.RegisterRunTimeManager(),
		CMImpl.RegisterDefaultConnectionManager(),
		provide.RegisterHttpHandler(settings.httpListenerUrl),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		ProvideTextListener(settings, ConsumerCounter),
		ProvideCompressedListener(settings, ConsumerCounter),
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
	)

	return NewTerminalAppUsingFxApp(fxApp, terminalApplication, shutDowner), nil
}
