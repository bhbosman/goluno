package lunoStream

import (
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/goLuno/internal/common"
)

type AppSettings struct {
	pairs                     []*common.PairInformation
	textListenerEnabled       bool
	textListenerUrl           string
	compressedListenerEnabled bool
	compressedListenerUrl     string
	canDial                   []goCommsNetDialer.ICanDial
	macConnections            int
	errors                    []error
}

type ILunoStreamAppApplySettings interface {
	apply(settings *AppSettings) error
}
