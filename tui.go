package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app     *tview.Application
	msgView *tview.TextView
)

func initTUI(onSend func(string)) {
	app = tview.NewApplication()

	msgView = tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true).
		ScrollToEnd()
	msgView.SetBorder(true)

	inputBox := tview.NewInputField().SetLabel("> ")

	inputBox.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}
		text := inputBox.GetText()
		if text == "" {
			return
		}
		inputBox.SetText("")
		go onSend(text)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(msgView, 0, 1, false).
		AddItem(inputBox, 1, 0, true)

	app.SetRoot(flex, true).SetFocus(inputBox)
}

func tuiPrint(line string) {
	app.QueueUpdateDraw(func() {
		fmt.Fprintf(msgView, "%s\n", line)
	})
}

func runTUI() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
