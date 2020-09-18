package fullMarketData

import (
	"github.com/emirpasic/gods/lists/singlylinkedlist"
)

type PricePoint struct {
	Price   float64
	List    *singlylinkedlist.List
	Touched bool
}

func NewPricePoint(price float64) *PricePoint {
	return &PricePoint{
		Price:   price,
		List:    singlylinkedlist.New(),
		Touched: false,
	}
}

func (self *PricePoint) AddOrder(order IOrder) {
	self.List.Add(order)
	self.Touched = true
}

func (self *PricePoint) Delete(id string) {
	index, _ := self.List.Find(
		func(index int, value interface{}) bool {
			if order, ok := value.(IOrder); ok {
				return order.GetId() == id
			}
			return false
		})
	if index >= 0 {
		self.List.Remove(index)
	}
	self.Touched = true
}

func (self *PricePoint) Count() int {
	return self.List.Size()
}

func (self *PricePoint) Find(id string) (bool, IOrder) {
	idx, value := self.List.Find(
		func(index int, value interface{}) bool {
			if order, ok := value.(IOrder); ok {
				return order.GetId() == id
			}
			return false
		})
	switch {
	case idx >= 0:
		order, ok := value.(IOrder)
		return ok, order
	default:
		return false, nil
	}
}

func (self *PricePoint) GetVolume() (pricePointVolume float64) {
	pricePointVolume = 0
	iterator := self.List.Iterator()
	for iterator.Next() {
		if order, ok := iterator.Value().(IOrder); ok {
			pricePointVolume += order.GetVolume()
		}
	}
	return pricePointVolume
}

func (self *PricePoint) ClearTouched() {
	self.Touched = false
}
