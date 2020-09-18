package common

import "fmt"

type PairInformation struct {
	Pair string
}

func NewPairInformation(pair string) *PairInformation {
	return &PairInformation{
		Pair: pair,
	}
}

func PublishName(s string) string {
	return fmt.Sprintf("Publish%v", s)
}

func RepublishName(s string) string {
	return fmt.Sprintf("Republish%v", s)
}
