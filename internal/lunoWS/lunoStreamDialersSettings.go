package lunoWS

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/gocomms/netDial"
)

type lunoStreamDialersSettings struct {
	pairs          []*common.PairInformation
	canDial        []netDial.ICanDial
	maxConnections int
}
