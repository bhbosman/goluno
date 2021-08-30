package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/common"
	"github.com/bhbosman/gocomms/netDial"
	"log"
)

type AppSettings struct {
	logger                    *log.Logger
	pairs                     []*common.PairInformation
	textListenerEnabled       bool
	textListenerUrl           string
	compressedListenerEnabled bool
	compressedListenerUrl     string
	httpListenerUrl           string
	canDial                   []netDial.ICanDial
	macConnections            int
}

type ILunoStreamAppApplySettings interface {
	apply(settings *AppSettings)
}
