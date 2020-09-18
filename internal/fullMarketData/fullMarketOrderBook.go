package fullMarketData

import (
	stream2 "github.com/bhbosman/goMessages/luno/stream"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/emirpasic/gods/utils"
)

type OrderSide int8

const BuySide OrderSide = 0
const AskSide OrderSide = 1

type FullMarketOrderBook struct {
	Orders map[string]struct {
		OrderSide
		*PricePoint
	}
	OrderSide          [2]*avltree.Tree
	SourceTimestamp    int64
	SourceMessageCount int64
}

func (self *FullMarketOrderBook) Clear() {
	self.Orders = make(map[string]struct {
		OrderSide
		*PricePoint
	})
	self.SourceTimestamp = 0
	self.OrderSide[0].Clear()
	self.OrderSide[1].Clear()
}

func (self *FullMarketOrderBook) AddOrder(side OrderSide, order IOrder) {
	get, found := self.OrderSide[uint8(side)].Get(order.GetPrice())
	if found {
		if pricePoint, ok := get.(*PricePoint); ok {
			pricePoint.AddOrder(order)
			self.Orders[order.GetId()] = struct {
				OrderSide
				*PricePoint
			}{side, pricePoint}
		}
	} else {
		pricePoint := NewPricePoint(order.GetPrice())
		pricePoint.AddOrder(order)
		self.OrderSide[uint8(side)].Put(order.GetPrice(), pricePoint)
		self.Orders[order.GetId()] = struct {
			OrderSide
			*PricePoint
		}{side, pricePoint}
	}
}

var epsilon = 1e-8

func (self *FullMarketOrderBook) TradeUpdate(tradeUpdate *stream2.TradeUpdate) {
	if makerOrder, ok := self.Orders[tradeUpdate.MakerOrderId]; ok {
		if find, order := makerOrder.PricePoint.Find(tradeUpdate.MakerOrderId); find {
			newVolume := order.ReduceVolume(tradeUpdate.Base)
			if newVolume <= epsilon {
				self.deleteOrder(tradeUpdate.MakerOrderId)
			}
		}
	}
}

func (self *FullMarketOrderBook) DeleteUpdate(update *stream2.DeleteUpdate) {
	self.deleteOrder(update.OrderId)
}

func (self *FullMarketOrderBook) deleteOrder(oderId string) {
	if data, ok := self.Orders[oderId]; ok {
		delete(self.Orders, oderId)
		data.PricePoint.Delete(oderId)
		if data.PricePoint.Count() == 0 {
			self.OrderSide[int8(data.OrderSide)].Remove(data.PricePoint.Price)
		}
	}
}

func (self *FullMarketOrderBook) CreateUpdate(update *stream2.CreateUpdate) {
	if update.Type == "BID" {
		self.AddOrder(BuySide, update)
		return
	}
	self.AddOrder(AskSide, update)
}

func (self *FullMarketOrderBook) SetTimeStamp(timestamp int64) {
	self.SourceTimestamp = timestamp
}

func (self *FullMarketOrderBook) UpdateMessageReceivedCount() {
	self.SourceMessageCount++
}

func NewFullMarketOrderBook() *FullMarketOrderBook {
	return &FullMarketOrderBook{
		Orders: make(map[string]struct {
			OrderSide
			*PricePoint
		}),
		OrderSide: [2]*avltree.Tree{
			avltree.NewWith(utils.Float64Comparator),
			avltree.NewWith(utils.Float64Comparator)},
	}
}
