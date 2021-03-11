package lunoStream

import (
	"github.com/bhbosman/goLuno/internal"
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/connectionManager/endpoints"
	"github.com/bhbosman/gocomms/connectionManager/view"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/gocomms/provide"
	"github.com/bhbosman/gologging"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"log"
)

func App(pairs ...ILunoStreamAppApplySettings) (*fx.App, fx.Shutdowner) {
	settings := &AppSettings{
		logger:                log.New(&stream.NullWriter{}, "", log.LstdFlags),
		pairs:                 nil,
		textListenerUrl:       "tcp4://127.0.0.1:3000",
		compressedListenerUrl: "tcp4://127.0.0.1:3001",
		httpListenerUrl:       "http://127.0.0.1:8080",
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
	fxApp := fx.New(
		fx.Supply(settings, ConsumerCounter),
		fx.Logger(settings.logger),
		gologging.ProvideLogFactory(settings.logger, nil),
		fx.Populate(&shutDowner),
		app2.RegisterRootContext(pubSub),
		connectionManager.RegisterDefaultConnectionManager(),
		provide.RegisterHttpHandler(settings.httpListenerUrl),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		impl.RegisterAllConnectionRelatedServices(),
		TextListener(pubSub, ConsumerCounter, 1024, settings.textListenerUrl, settings.pairs...),
		CompressedListener(pubSub, ConsumerCounter, 1024, settings.compressedListenerUrl, settings.pairs...),
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
