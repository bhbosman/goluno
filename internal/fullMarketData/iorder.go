package fullMarketData

type IOrder interface {
	GetId() string
	GetPrice() float64
	GetVolume() float64
	ReduceVolume(base float64) (leftOverVolume float64)
}
