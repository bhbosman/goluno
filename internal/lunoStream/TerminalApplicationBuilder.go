package lunoStream

import "github.com/rivo/tview"

type TerminalApplicationBuilder struct {
}

func (self *TerminalApplicationBuilder) createCommandList() (commandList *tview.List) {
	///// Commands /////
	commandList = tview.NewList()
	commandList.SetBorder(true).SetTitle("Commands")
	commandList.ShowSecondaryText(false)
	return commandList
}

func (self *TerminalApplicationBuilder) createOutputPanel(app *tview.Application) (outputPanel *tview.Flex) {
	outputPanel = tview.NewFlex().SetDirection(tview.FlexRow)
	return outputPanel
}

func (self *TerminalApplicationBuilder) createMainLayout(
	commandList tview.Primitive,
	outputPanel tview.Primitive,
	headers ...tview.Primitive) (layout *tview.Flex) {
	///// Main Layout /////
	mainLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(commandList, 30, 1, true).
		AddItem(outputPanel, 0, 4, false)

	info := tview.NewTextView()
	info.SetText("Terminal app")
	//info.SetTextAlign(tview.AlignRight)

	//info2 := tview.NewTextView()
	//info2.SetText("Time: ")

	layout2 := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout2.AddItem(info, 0, 1, false)
	for _, header := range headers {
		layout2.AddItem(header, 0, 1, false)
	}

	layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(layout2, 1, 1, false).
		AddItem(mainLayout, 0, 20, true)

	return layout
}
