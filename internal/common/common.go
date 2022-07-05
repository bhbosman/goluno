package common

import (
	"fmt"
	"github.com/bhbosman/gocommon/model"
	"net/url"
)

type PairInformation struct {
	UseSocks5          bool
	SocksUrl           *url.URL
	PairUrl            *url.URL
	Pair               string
	ServiceIdentifier  model.ServiceIdentifier
	ServiceDependentOn model.ServiceIdentifier
}

func NewPairInformation(
	pair string,
	ServiceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	UseSocks5 bool,
	SocksUrl string,
	PairUrl string,
) (*PairInformation, error) {
	sockUrl, err := url.Parse(SocksUrl)
	if err != nil {
		return nil, err
	}
	pairUrl, err := url.Parse(PairUrl)
	if err != nil {
		return nil, err
	}
	return &PairInformation{
		UseSocks5:          UseSocks5,
		SocksUrl:           sockUrl,
		PairUrl:            pairUrl,
		Pair:               pair,
		ServiceIdentifier:  ServiceIdentifier,
		ServiceDependentOn: serviceDependentOn,
	}, nil
}

func PublishName(s string) string {
	return fmt.Sprintf("Publish%v", s)
}

func RepublishName(s string) string {
	return fmt.Sprintf("Republish%v", s)
}
