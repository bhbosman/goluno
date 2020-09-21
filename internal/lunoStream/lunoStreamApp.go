package lunoStream

import (
	"github.com/bhbosman/goLuno/internal"
	"github.com/bhbosman/goLuno/internal/ConsumerCounter"
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/comms/connectionManager/endpoints"
	"github.com/bhbosman/gocommon/comms/connectionManager/view"
	"github.com/bhbosman/gocommon/comms/http"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
	"log"
	"os"
)

func App(pairs ...ILunoStreamAppApplySettings) (*fx.App, fx.Shutdowner) {
	settings := &AppSettings{
		logger:                log.New(os.Stdout, "", log.LstdFlags),
		pairs:                 nil,
		textListenerUrl:       "tcp4://127.0.0.1:3000",
		compressedListenerUrl: "tcp4://127.0.0.1:3001",
		httpListenerUrl:       "http://127.0.0.1:8080",
	}
	for _, apply := range pairs {
		apply.apply(settings)
	}
	ConsumerCounter := &ConsumerCounter.ConsumerCounter{}
	var shutDowner fx.Shutdowner
	fxApp := fx.New(
		fx.Supply(settings, ConsumerCounter),
		fx.Logger(settings.logger),
		gologging.ProvideLogFactory(settings.logger, nil),
		fx.Populate(&shutDowner),
		app2.RegisterRootContext(),
		connectionManager.RegisterDefaultConnectionManager(),
		http.RegisterHttpHandler(settings.httpListenerUrl),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		commsImpl.RegisterAllConnectionRelatedServices(),
		TextListener(settings.textListenerUrl, settings.pairs...),
		CompressedListener(settings.compressedListenerUrl, settings.pairs...),
		Dialers(CanDial(ConsumerCounter), AddPairsInformation(settings.pairs)),
		ProvideReadLunoKeys(),
		internal.InvokeApps(),
	)
	return fxApp, shutDowner
}
