package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/common"
	"log"
)

type addCurrencyPair struct {
	pair *common.PairInformation
}

func AddCurrencyPair(pair *common.PairInformation) *addCurrencyPair {
	return &addCurrencyPair{pair: pair}
}

func (self addCurrencyPair) apply(settings *AppSettings) {
	settings.pairs = append(settings.pairs, self.pair)
}

type textListenerUrl struct {
	url string
}

func TextListenerUrl(url string) *textListenerUrl {
	return &textListenerUrl{url: url}
}

func (self textListenerUrl) apply(settings *AppSettings) {
	settings.textListenerUrl = self.url
	settings.textListenerEnabled = true
}

type compressedListenerUrl struct {
	url string
}

func CompressedListenerUrl(url string) *compressedListenerUrl {
	return &compressedListenerUrl{url: url}
}

func (self compressedListenerUrl) apply(settings *AppSettings) {
	settings.compressedListenerUrl = self.url
	settings.compressedListenerEnabled = true
}

type httpListenerUrl struct {
	url string
}

func HttpListenerUrl(url string) *httpListenerUrl {
	return &httpListenerUrl{url: url}
}

func (self httpListenerUrl) apply(settings *AppSettings) {
	settings.httpListenerUrl = self.url
}

type setlogger struct {
	logger *log.Logger
}

func (self setlogger) apply(settings *AppSettings) {
	settings.logger = self.logger
}

func Logger(logger *log.Logger) *setlogger {
	return &setlogger{logger: logger}
}
