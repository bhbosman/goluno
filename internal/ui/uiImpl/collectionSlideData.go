package uiImpl

import (
	"context"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/connectionManager/model"
	"sort"
	"time"
)

type PublishInstanceDataFor struct {
	Id   string
	Name string
}

type IdAndName struct {
	Id   string
	Name string
}

type iConnectionData interface {
	ISendMessage.ISendMessage
}
type ConnectionData struct {
	ConnectionId   string
	isDirty        bool
	CancelContext  context.Context
	CancelFunc     context.CancelFunc
	Name           string
	Status         string
	ConnectionTime time.Time
	Grid           []model.LineData
}

type collectionSlideData struct {
	connectionListIsDirty      bool
	ss                         map[string]*ConnectionData
	messageRouter              *messageRouter.MessageRouter
	onConnectionListChange     func(connectionList []IdAndName)
	onConnectionInstanceChange func(data *ConnectionData)
}

func NewData() *collectionSlideData {
	result := &collectionSlideData{
		ss:            make(map[string]*ConnectionData),
		messageRouter: messageRouter.NewMessageRouter(),
	}
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	_ = result.messageRouter.Add(result.handleConnectionState)
	_ = result.messageRouter.Add(result.handleConnectionCreated)
	_ = result.messageRouter.Add(result.handleConnectionClosed)
	_ = result.messageRouter.Add(result.handlePublishInstanceDataFor)
	return result
}

func (self *collectionSlideData) Send(data interface{}) error {
	_, err := self.messageRouter.Route(data)
	return err
}

func (self *collectionSlideData) handlePublishInstanceDataFor(message *PublishInstanceDataFor) error {
	if info, ok := self.ss[message.Id]; ok {
		self.DoConnectionInstanceChange(info)
	}
	return nil
}
func (self *collectionSlideData) handleEmptyQueue(data *messages.EmptyQueue) error {
	if self.connectionListIsDirty {
		self.DoConnectionListChange()
		self.connectionListIsDirty = false
	}
	for _, connectionData := range self.ss {
		if connectionData.isDirty {
			self.DoConnectionInstanceChange(connectionData)
			connectionData.isDirty = false
		}
	}
	return nil
}

func (self *collectionSlideData) handleConnectionState(message *model.ConnectionState) error {
	if data, ok := self.ss[message.ConnectionId]; ok {
		data.isDirty = true
		data.CancelContext = message.CancelContext
		data.CancelFunc = message.CancelFunc
		data.Name = message.Name
		data.Status = message.Status
		data.ConnectionTime = message.ConnectionTime
		data.Grid = message.Grid
	}
	return nil
}

func (self *collectionSlideData) handleConnectionClosed(message *model.ConnectionClosed) error {
	delete(self.ss, message.ConnectionId)
	self.connectionListIsDirty = true
	return nil
}

func (self *collectionSlideData) handleConnectionCreated(message *model.ConnectionCreated) error {
	self.ss[message.ConnectionId] = &ConnectionData{
		isDirty:      true,
		ConnectionId: message.ConnectionId}
	self.connectionListIsDirty = true
	return nil
}

func (self *collectionSlideData) DoConnectionListChange() {
	if self.onConnectionListChange != nil {
		ss := make([]string, 0, len(self.ss))

		for key, _ := range self.ss {
			ss = append(ss, key)
		}
		sort.Strings(ss)
		cbData := make([]IdAndName, 0, len(self.ss))
		for key, data := range self.ss {
			idAndName := IdAndName{
				Id:   key,
				Name: data.Name,
			}
			cbData = append(cbData, idAndName)
		}
		self.onConnectionListChange(cbData)
	}
}
func (self *collectionSlideData) DoConnectionInstanceChange(data *ConnectionData) {
	if self.onConnectionInstanceChange != nil {
		self.onConnectionInstanceChange(data)
	}
}

func (self *collectionSlideData) SetConnectionInstanceChange(cb func(data *ConnectionData)) {
	self.onConnectionInstanceChange = cb
}
func (self *collectionSlideData) SetConnectionListChange(cb func(connectionList []IdAndName)) {
	self.onConnectionListChange = cb
}
