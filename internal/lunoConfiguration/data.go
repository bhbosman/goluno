package lunoConfiguration

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
)

type data struct {
	//m             map[string]*LunoConfiguration
	isDirty       map[string]bool
	MessageRouter messageRouter.IMessageRouter
}

func (self *data) Send(message interface{}) error {
	self.MessageRouter.Route(message)
	return nil
}

func (self *data) SomeMethod() {
}

func (self *data) ShutDown() error {
	return nil
}

func (self *data) handleEmptyQueue(*messages.EmptyQueue) {
	self.isDirty = make(map[string]bool)
}

func newData() (ILunoConfigurationData, error) {
	result := &data{
		MessageRouter: messageRouter.NewMessageRouter(),
		isDirty:       make(map[string]bool),
	}
	_ = result.MessageRouter.Add(result.handleEmptyQueue)
	return result, nil
}
