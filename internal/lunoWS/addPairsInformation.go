package lunoWS

import "github.com/bhbosman/goLuno/internal/common"

type addPairsInformation struct {
	pairs []*common.PairInformation
}

func AddPairsInformation(pairs []*common.PairInformation) *addPairsInformation {
	return &addPairsInformation{pairs: pairs}
}

func (self addPairsInformation) apply(settings *lunoStreamDialersSettings) {
	for _, pair := range self.pairs {
		settings.pairs = append(settings.pairs, pair)
	}
}
