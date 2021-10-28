package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/common"
)

type addCurrencyPair struct {
	pair *common.PairInformation
}

func AddCurrencyPair(pair *common.PairInformation) *addCurrencyPair {
	return &addCurrencyPair{pair: pair}
}

func (self *addCurrencyPair) apply(settings *AppSettings) error {
	settings.pairs = append(settings.pairs, self.pair)
	return nil
}

type textListenerUrl struct {
	url string
}

func TextListenerUrl(url string) *textListenerUrl {
	return &textListenerUrl{url: url}
}

func (self *textListenerUrl) apply(settings *AppSettings) error {
	settings.textListenerUrl = self.url
	settings.textListenerEnabled = true
	return nil
}

type compressedListenerUrl struct {
	url string
}

func CompressedListenerUrl(url string) *compressedListenerUrl {
	return &compressedListenerUrl{url: url}
}

func (self *compressedListenerUrl) apply(settings *AppSettings) error {
	settings.compressedListenerUrl = self.url
	settings.compressedListenerEnabled = true
	return nil
}

type httpListenerUrl struct {
	url string
}

func HttpListenerUrl(url string) *httpListenerUrl {
	return &httpListenerUrl{url: url}
}

func (self *httpListenerUrl) apply(settings *AppSettings) error {
	settings.httpListenerUrlEnabled = true
	settings.httpListenerUrl = self.url
	return nil
}
