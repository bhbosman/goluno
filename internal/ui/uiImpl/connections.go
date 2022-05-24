package uiImpl

import (
	"github.com/rivo/tview"
)

func (self *Service) Connections() Slide {
	return func(nextSlide func()) (title string, content tview.Primitive) {
		return "Connection", NewConnectionSlide(self.ctx, self.pubSub, self.app)

	}
}
