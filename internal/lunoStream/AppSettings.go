package lunoStream

import (
	"github.com/bhbosman/goCommsNetDialer"
)

type AppSettings struct {
	//textListenerEnabled       bool
	//textListenerUrl           string
	compressedListenerEnabled bool
	compressedListenerUrl     string
	canDial                   []goCommsNetDialer.ICanDial
	macConnections            int
	errors                    []error
}

type ILunoStreamAppApplySettings interface {
	apply(settings *AppSettings) error
}
