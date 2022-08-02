// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/goLuno/internal/lunoConfiguration (interfaces: ILunoConfiguration)

// Package lunoConfiguration is a generated GoMock package.
package lunoConfiguration

import (
	fmt "fmt"

	errors "github.com/bhbosman/gocommon/errors"
	"golang.org/x/net/context"
)

// Interface A Comment
// Interface github.com/bhbosman/goLuno/internal/lunoConfiguration
// Interface ILunoConfiguration
// Interface ILunoConfiguration, Method: GetAll
type ILunoConfigurationGetAllIn struct {
}

type ILunoConfigurationGetAllOut struct {
	//Args0 []*LunoConfiguration
}
type ILunoConfigurationGetAllError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *ILunoConfigurationGetAllError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type ILunoConfigurationGetAll struct {
	inData         ILunoConfigurationGetAllIn
	outDataChannel chan ILunoConfigurationGetAllOut
}

func NewILunoConfigurationGetAll(waitToComplete bool) *ILunoConfigurationGetAll {
	var outDataChannel chan ILunoConfigurationGetAllOut
	if waitToComplete {
		outDataChannel = make(chan ILunoConfigurationGetAllOut)
	} else {
		outDataChannel = nil
	}
	return &ILunoConfigurationGetAll{
		inData:         ILunoConfigurationGetAllIn{},
		outDataChannel: outDataChannel,
	}
}

func (self *ILunoConfigurationGetAll) Wait(onError func(interfaceName string, methodName string, err error) error) (ILunoConfigurationGetAllOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &ILunoConfigurationGetAllError{
			InterfaceName: "ILunoConfiguration",
			MethodName:    "GetAll",
			Reason:        "Channel for ILunoConfiguration::GetAll returned false",
		}
		if onError != nil {
			err := onError("ILunoConfiguration", "GetAll", generatedError)
			return ILunoConfigurationGetAllOut{}, err
		} else {
			return ILunoConfigurationGetAllOut{}, generatedError
		}
	}
	return data, nil
}

func (self *ILunoConfigurationGetAll) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallILunoConfigurationGetAll(context context.Context, channel chan<- interface{}, waitToComplete bool) (ILunoConfigurationGetAllOut, error) {
	if context != nil && context.Err() != nil {
		return ILunoConfigurationGetAllOut{}, context.Err()
	}
	data := NewILunoConfigurationGetAll(waitToComplete)
	if waitToComplete {
		defer func(data *ILunoConfigurationGetAll) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return ILunoConfigurationGetAllOut{}, context.Err()
	}
	channel <- data
	var err error
	var v ILunoConfigurationGetAllOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return ILunoConfigurationGetAllOut{}, err
	}
	return v, nil
}

// Interface ILunoConfiguration, Method: Send
type ILunoConfigurationSendIn struct {
	arg0 interface{}
}

type ILunoConfigurationSendOut struct {
	Args0 error
}
type ILunoConfigurationSendError struct {
	InterfaceName string
	MethodName    string
	Reason        string
}

func (self *ILunoConfigurationSendError) Error() string {
	return fmt.Sprintf("error in data coming back from %v::%v. Reason: %v", self.InterfaceName, self.MethodName, self.Reason)
}

type ILunoConfigurationSend struct {
	inData         ILunoConfigurationSendIn
	outDataChannel chan ILunoConfigurationSendOut
}

func NewILunoConfigurationSend(waitToComplete bool, arg0 interface{}) *ILunoConfigurationSend {
	var outDataChannel chan ILunoConfigurationSendOut
	if waitToComplete {
		outDataChannel = make(chan ILunoConfigurationSendOut)
	} else {
		outDataChannel = nil
	}
	return &ILunoConfigurationSend{
		inData: ILunoConfigurationSendIn{
			arg0: arg0,
		},
		outDataChannel: outDataChannel,
	}
}

func (self *ILunoConfigurationSend) Wait(onError func(interfaceName string, methodName string, err error) error) (ILunoConfigurationSendOut, error) {
	data, ok := <-self.outDataChannel
	if !ok {
		generatedError := &ILunoConfigurationSendError{
			InterfaceName: "ILunoConfiguration",
			MethodName:    "Send",
			Reason:        "Channel for ILunoConfiguration::Send returned false",
		}
		if onError != nil {
			err := onError("ILunoConfiguration", "Send", generatedError)
			return ILunoConfigurationSendOut{}, err
		} else {
			return ILunoConfigurationSendOut{}, generatedError
		}
	}
	return data, nil
}

func (self *ILunoConfigurationSend) Close() error {
	close(self.outDataChannel)
	return nil
}
func CallILunoConfigurationSend(context context.Context, channel chan<- interface{}, waitToComplete bool, arg0 interface{}) (ILunoConfigurationSendOut, error) {
	if context != nil && context.Err() != nil {
		return ILunoConfigurationSendOut{}, context.Err()
	}
	data := NewILunoConfigurationSend(waitToComplete, arg0)
	if waitToComplete {
		defer func(data *ILunoConfigurationSend) {
			err := data.Close()
			if err != nil {
			}
		}(data)
	}
	if context != nil && context.Err() != nil {
		return ILunoConfigurationSendOut{}, context.Err()
	}
	channel <- data
	var err error
	var v ILunoConfigurationSendOut
	if waitToComplete {
		v, err = data.Wait(func(interfaceName string, methodName string, err error) error {
			return err
		})
	} else {
		err = errors.NoWaitOperationError
	}
	if err != nil {
		return ILunoConfigurationSendOut{}, err
	}
	return v, nil
}

func ChannelEventsForILunoConfiguration(next ILunoConfiguration, event interface{}) (bool, error) {
	switch v := event.(type) {
	case *ILunoConfigurationSend:
		data := ILunoConfigurationSendOut{}
		data.Args0 = next.Send(v.inData.arg0)
		if v.outDataChannel != nil {
			v.outDataChannel <- data
		}
		return true, nil
	default:
		return false, nil
	}
}