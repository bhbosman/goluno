package lunoWS

import (
	"github.com/bhbosman/goCommsNetDialer"
)

type canDialSetting struct {
	canDial []goCommsNetDialer.ICanDial
}

func CanDial(canDial ...goCommsNetDialer.ICanDial) *canDialSetting {
	return &canDialSetting{canDial: canDial}
}

func (self canDialSetting) apply(settings *lunoStreamDialersSettings) {
	for _, cd := range self.canDial {
		settings.canDial = append(settings.canDial, cd)
	}
}
