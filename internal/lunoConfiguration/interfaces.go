package lunoConfiguration

import (
	"github.com/bhbosman/gocommon/services/IDataShutDown"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

//type LunoConfiguration struct {
//	Name string
//}

//func NewLunoConfiguration(name string) *LunoConfiguration {
//	return &LunoConfiguration{Name: name}
//}

type ILunoConfiguration interface {
	ISendMessage.ISendMessage
	//GetAll() []*LunoConfiguration
}

type ILunoConfigurationService interface {
	ILunoConfiguration
	IFxService.IFxServices
}

type ILunoConfigurationData interface {
	ILunoConfiguration
	IDataShutDown.IDataShutDown
	ISendMessage.ISendMessage
}
