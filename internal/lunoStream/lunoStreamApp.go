package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/lunoWS"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/Services/implementations"
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/connectionManager/endpoints"
	"github.com/bhbosman/gocomms/connectionManager/view"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/provide"
	"go.uber.org/fx"
	"go.uber.org/multierr"
)

func ProvideTextListener(settings *AppSettings, ConsumerCounter *netDial.CanDialDefaultImpl) fx.Option {
	var opts []fx.Option
	if settings.textListenerEnabled {
		opts = append(opts, TextListener(ConsumerCounter, 1024, settings.textListenerUrl, settings.pairs...))
	}
	return fx.Options(opts...)
}

func ProvideCompressedListener(settings *AppSettings, ConsumerCounter *netDial.CanDialDefaultImpl) fx.Option {
	var opts []fx.Option
	if settings.compressedListenerEnabled {
		opts = append(opts,
			CompressedListener(
				ConsumerCounter,
				1024,
				settings.compressedListenerUrl,
				settings.pairs...))
	}
	return fx.Options(opts...)
}

func App(pairs ...ILunoStreamAppApplySettings) (*LunaApp, error) {
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

	fxApp := fx.New(
		fx.Supply(settings, ConsumerCounter),
		fx.Populate(&shutDowner),
		fx.Populate(&runTimeManager),

		app2.ProvideZapCoreEncoderConfigForDev(),
		app2.ProvideZapCoreEncoderConfigForProd(),
		app2.ProvideZapConfigForDev(),
		app2.ProvideZapConfigForProd(),
		app2.ProvideZapLogger(),
		app2.ProvideFxWithLogger(),
		app2.RegisterRunTimeManager(),
		app2.RegisterApplicationContext(),
		app2.ProvidePubSub("Application"),

		connectionManager.RegisterDefaultConnectionManager(),
		provide.RegisterHttpHandler(settings.httpListenerUrl),
		endpoints.RegisterConnectionManagerEndpoint(),
		implementations.ProvideNewUniqueReferenceService(),
		implementations.ProvideUniqueSessionNumber(),
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
		app2.InvokeApps(),
	)
	return &LunaApp{
		FxApp:      fxApp,
		ShutDowner: shutDowner,
	}, nil
}
