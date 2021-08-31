package lunoStream

import (
	"github.com/bhbosman/goLuno/internal"
	"github.com/bhbosman/gocommon"
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/connectionManager/endpoints"
	"github.com/bhbosman/gocomms/connectionManager/view"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/provide"
	"github.com/bhbosman/gologging"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"log"
	"os"
)

func ProvideTextListener(settings *AppSettings, pubSub *pubsub.PubSub, ConsumerCounter *netDial.CanDialDefaultImpl) fx.Option {
	var opts []fx.Option
	if settings.textListenerEnabled {
		opts = append(opts, TextListener(pubSub, ConsumerCounter, 1024, settings.textListenerUrl, settings.pairs...))
	}
	return fx.Options(opts...)
}

func ProvideCompressedListener(settings *AppSettings, pubSub *pubsub.PubSub, ConsumerCounter *netDial.CanDialDefaultImpl) fx.Option {
	var opts []fx.Option
	if settings.compressedListenerEnabled {
		opts = append(opts,
			CompressedListener(
				pubSub,
				ConsumerCounter,
				1024,
				settings.compressedListenerUrl,
				settings.pairs...))
	}
	return fx.Options(opts...)
}

func App(pairs ...ILunoStreamAppApplySettings) (*fx.App, fx.Shutdowner) {
	settings := &AppSettings{
		//logger:                log.New(&stream.NullWriter{}, "", log.LstdFlags),
		logger:                    log.New(os.Stderr, "", log.LstdFlags),
		pairs:                     nil,
		textListenerEnabled:       false,
		textListenerUrl:           "tcp4://127.0.0.1:3000",
		compressedListenerEnabled: false,
		compressedListenerUrl:     "tcp4://127.0.0.1:3001",
		httpListenerUrl:           "http://127.0.0.1:8080",
		canDial:                   nil,
		macConnections:            0,
	}
	for _, apply := range pairs {
		apply.apply(settings)
	}
	ConsumerCounter := netDial.NewCanDialDefaultImpl()
	pubSub := pubsub.New(32)
	lunoKeys, err := ReadLunoKeys()
	if err != nil {
		return fx.New(fx.Error(err)), nil
	}
	var shutDowner fx.Shutdowner
	var dd *gocommon.RunTimeManager

	fxApp := fx.New(
		fx.Supply(settings, ConsumerCounter),
		fx.Logger(settings.logger),
		gologging.ProvideLogFactory(settings.logger, nil),
		fx.Populate(&shutDowner),
		fx.Populate(&dd),
		app2.RegisterRootContext(pubSub),
		connectionManager.RegisterDefaultConnectionManager(),
		provide.RegisterHttpHandler(settings.httpListenerUrl),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		ProvideTextListener(settings, pubSub, ConsumerCounter),
		ProvideCompressedListener(settings, pubSub, ConsumerCounter),

		Dialers(
			lunoKeys.Key,
			lunoKeys.Secret,
			pubSub,
			ConsumerCounter,
			CanDial(ConsumerCounter),
			AddPairsInformation(settings.pairs),
			MaxConnections(1)),
		ProvideReadLunoKeys(lunoKeys),
		internal.InvokeApps(),
	)
	return fxApp, shutDowner
}
