package lunoWS

import (
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goLuno/internal/common"
)

type lunoStreamDialersSettings struct {
	pairs          []*common.PairInformation
	canDial        []goCommsNetDialer.ICanDial
	maxConnections int
}
