package uiIntf

import (
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/rivo/tview"
)

type OnApplication func() *tview.Application
type IUi interface {
	Build() OnApplication
}

type IUiService interface {
	IUi
	//interfaces.IFxServices
}

type IUiData interface {
	IUi
	interfaces.IDataShutDown
	interfaces.ISendMessage
}
