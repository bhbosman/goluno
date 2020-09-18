package internal

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/goLuno/internal/provide"
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/comms/connectionManager/endpoints"
	"github.com/bhbosman/gocommon/comms/connectionManager/view"
	"github.com/bhbosman/gocommon/comms/http"
	log2 "github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
	"log"
	"os"
)

type LunoStreamAppSettings struct {
	logger                *log.Logger
	pairs                 []*common.PairInformation
	textListenerUrl       string
	compressedListenerUrl string
	httpListenerUrl       string
}

func LunoStreamApp(pairs ...ILunoStreamAppApplySettings) (*fx.App, fx.Shutdowner) {
	settings := &LunoStreamAppSettings{
		logger:                log.New(os.Stdout, "", log.LstdFlags),
		pairs:                 nil,
		textListenerUrl:       "tcp4://127.0.0.1:3000",
		compressedListenerUrl: "tcp4://127.0.0.1:3001",
		httpListenerUrl:       "http://127.0.0.1:8080",
	}
	for _, apply := range pairs {
		apply.apply(settings)
	}
	var shutDowner fx.Shutdowner
	fxApp := fx.New(
		fx.Logger(settings.logger),
		log2.ProvideLogFactory(settings.logger, nil),
		//fx.LogName("Luno Stream Application"),
		fx.Populate(&shutDowner),
		app2.RegisterRootContext(),
		connectionManager.RegisterDefaultConnectionManager(),
		http.RegisterHttpHandler(settings.httpListenerUrl),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		commsImpl.RegisterAllConnectionRelatedServices(),
		provide.TextListener(settings.textListenerUrl, settings.pairs...),
		provide.CompressedListener(settings.compressedListenerUrl, settings.pairs...),
		provide.LunoStreamDialers(settings.pairs...),
		provide.ProvideReadLunoKeys(),
		InvokeApps(),
	)
	return fxApp, shutDowner
}

type ILunoStreamAppApplySettings interface {
	apply(settings *LunoStreamAppSettings)
}
type addCurrencyPair struct {
	pair *common.PairInformation
}

func AddCurrencyPair(pair *common.PairInformation) *addCurrencyPair {
	return &addCurrencyPair{pair: pair}
}

func (self addCurrencyPair) apply(settings *LunoStreamAppSettings) {
	settings.pairs = append(settings.pairs, self.pair)
}

type textListenerUrl struct {
	url string
}

func TextListenerUrl(url string) *textListenerUrl {
	return &textListenerUrl{url: url}
}

func (self textListenerUrl) apply(settings *LunoStreamAppSettings) {
	settings.textListenerUrl = self.url
}

type compressedListenerUrl struct {
	url string
}

func CompressedListenerUrl(url string) *compressedListenerUrl {
	return &compressedListenerUrl{url: url}
}

func (self compressedListenerUrl) apply(settings *LunoStreamAppSettings) {
	settings.compressedListenerUrl = self.url
}

type httpListenerUrl struct {
	url string
}

func HttpListenerUrl(url string) *httpListenerUrl {
	return &httpListenerUrl{url: url}
}

func (self httpListenerUrl) apply(settings *LunoStreamAppSettings) {
	settings.httpListenerUrl = self.url
}

type setlogger struct {
	logger *log.Logger
}

func (self setlogger) apply(settings *LunoStreamAppSettings) {
	settings.logger = self.logger
}

func Logger(logger *log.Logger) *setlogger {
	return &setlogger{logger: logger}
}
