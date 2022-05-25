package uiImpl

import (
	"github.com/cskr/pubsub"
	"github.com/rivo/tview"
	"golang.org/x/net/context"
)

func Connections(applicationContext context.Context, pubSub *pubsub.PubSub, app *tview.Application) Slide {
	return func(nextSlide func()) (title string, content tview.Primitive) {
		return "Connection", NewConnectionSlide(applicationContext, pubSub, app)

	}
}
