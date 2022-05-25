package uiImpl

import (
	"context"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ConnectionSlide struct {
	data           iConnectionData
	connectionList *tview.List
	table          *tview.Table
	actionList     *tview.List
	next           tview.Primitive
	ctx            context.Context
	cancelFunc     context.CancelFunc
	channel        chan interface{}
	pubSub         *pubsub.PubSub
	app            *tview.Application
}

func (self *ConnectionSlide) Close() error {
	self.cancelFunc()
	close(self.channel)
	return nil
}

func (self *ConnectionSlide) Draw(screen tcell.Screen) {
	self.next.Draw(screen)
}

func (self *ConnectionSlide) GetRect() (int, int, int, int) {
	return self.next.GetRect()
}

func (self *ConnectionSlide) SetRect(x, y, width, height int) {
	self.next.SetRect(x, y, width, height)
}

func (self *ConnectionSlide) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return self.next.InputHandler()
}

func (self *ConnectionSlide) Focus(delegate func(p tview.Primitive)) {
	self.next.Focus(delegate)
}

func (self *ConnectionSlide) HasFocus() bool {
	return self.next.HasFocus()
}

func (self *ConnectionSlide) Blur() {
	self.next.Blur()
}

func (self *ConnectionSlide) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return self.next.MouseHandler()
}

func (self *ConnectionSlide) goRun() {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.channel)

	self.pubSub.AddSub(self.channel, "ActiveConnectionStatus")
loop:
	for {
		select {
		case <-self.ctx.Done():
			break loop
		case data, ok := <-self.channel:
			if !ok {
				break loop
			}
			success, _ := ISendMessage.ChannelEventsForISendMessage(self.data, data)
			if success {
				continue
			}
			_ = self.data.Send(data)

			if self.ctx.Err() != nil {
				return
			}
			if len(self.channel) == 0 {
				_ = self.data.Send(&messages.EmptyQueue{})
			}
		}
	}
	go func() {
		self.pubSub.Unsub(self.channel)
	}()
}

func (self *ConnectionSlide) SetConnectionListChange(list []IdAndName) {
	self.app.QueueUpdateDraw(func() {
		self.connectionList.Clear()
		for _, s := range list {
			self.connectionList.AddItem(s.Id, s.Name, 0, func() {
				self.app.SetFocus(self.actionList)
				self.actionList.SetDoneFunc(func() {
					self.app.SetFocus(self.connectionList)
				}).SetSelectedFunc(func(i int, s string, s2 string, r rune) {
					self.app.SetFocus(self.connectionList)
				})

			})
		}
	})
}

func (self *ConnectionSlide) SetConnectionInstanceChange(data *ConnectionData) {
	self.app.QueueUpdateDraw(func() {
		index := self.connectionList.GetCurrentItem()
		text, _ := self.connectionList.GetItemText(index)
		if text == data.ConnectionId {
			if data != nil && data.Grid != nil {
				tableData := newConnectionPlateContent(data.Grid)
				if tableData != nil {
					self.table.SetContent(tableData)
				}
			}
		}
	})
}

func NewConnectionSlide(applicationContext context.Context, pubSub *pubsub.PubSub, app *tview.Application) *ConnectionSlide {
	ctx, cancelFunc := context.WithCancel(applicationContext)
	channel := make(chan interface{}, 32)

	connectionList := tview.NewList().ShowSecondaryText(true)
	connectionList.SetBorder(true).SetTitle("Active Connections")
	connectionList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		_, _ = ISendMessage.CallISendMessageSend(ctx, channel, false, &PublishInstanceDataFor{
			Id:   mainText,
			Name: secondaryText,
		})
	})

	actionList := tview.NewList().ShowSecondaryText(false)
	actionList.SetBorder(true).SetTitle("Acions")
	actionList.AddItem("Disconnect", "", 0, func() {

	})
	table := tview.NewTable()
	table.SetBorder(true)
	layout := tview.NewFlex().
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(connectionList, 0, 3, true).
					AddItem(actionList, 4, 2, false),
					0,
					1,
					true).
				AddItem(table, 0, 4, false),
			0,
			1,
			true)
	data := NewData()
	result := &ConnectionSlide{
		data:           data,
		connectionList: connectionList,
		table:          table,
		actionList:     actionList,
		next:           layout,
		ctx:            ctx,
		cancelFunc:     cancelFunc,
		channel:        channel,
		pubSub:         pubSub,
		app:            app,
	}
	data.SetConnectionListChange(result.SetConnectionListChange)
	data.SetConnectionInstanceChange(result.SetConnectionInstanceChange)
	go result.goRun()
	return result
}
