package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/gocomms/netDial"
)

type AppSettings struct {
	pairs                     []*common.PairInformation
	textListenerEnabled       bool
	textListenerUrl           string
	compressedListenerEnabled bool
	compressedListenerUrl     string
	httpListenerUrlEnabled    bool
	httpListenerUrl           string
	canDial                   []netDial.ICanDial
	macConnections            int
}

type ILunoStreamAppApplySettings interface {
	apply(settings *AppSettings) error
}
