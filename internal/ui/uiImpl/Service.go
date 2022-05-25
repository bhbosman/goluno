package uiImpl

import (
	"context"
	"fmt"
	"github.com/bhbosman/goLuno/internal/ui/uiIntf"
	"github.com/cskr/pubsub"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type Slide func(nextSlide func()) (title string, content tview.Primitive)

type Service struct {
	app        *tview.Application
	ctx        context.Context
	cancelFunc context.CancelFunc
	pubSub     *pubsub.PubSub
}

func (self *Service) Build() uiIntf.OnApplication {
	return func() *tview.Application {
		build := self.BuildApp()
		return build
	}
}

func (self *Service) BuildApp() *tview.Application {
	slides := []Slide{
		self.Cover,
		Connections(self.ctx, self.pubSub, self.app),
	}

	pages := tview.NewPages()

	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			pages.SwitchToPage(added[0])
		})

	nextSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide + 1) % len(slides)
		info.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}
	previousSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide - 1 + len(slides)) % len(slides)
		info.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}

	for index, slide := range slides {
		title, primitive := slide(nextSlide)
		pages.AddPage(strconv.Itoa(index), primitive, true, index == 0)
		fmt.Fprintf(info, `%d ["%d"][green]%s[white][""]  `, index+1, index, title)
	}
	info.Highlight("0")

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	self.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlO {
			nextSlide()
			return nil
		} else if event.Key() == tcell.KeyCtrlP {
			previousSlide()
			return nil
		}
		return event
	})
	return self.app.SetRoot(layout, true).EnableMouse(true)
}

func NewService(
	ctx context.Context,
	pubSub *pubsub.PubSub) *Service {
	result := &Service{
		pubSub: pubSub,
		app:    tview.NewApplication(),
	}
	result.ctx, result.cancelFunc = context.WithCancel(ctx)
	return result
}
