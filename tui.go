package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app     *tview.Application
	msgView map[string]*tview.TextView
	listNum int = 5
	activeChannel string = "channel 0"
	currentUser string//both of these are assigned in sign in window
	currentPass string//passed tp client.go
)

func initTUI(onSend func(string)) {
	app = tview.NewApplication()
	pages := tview.NewPages()
	msgView = make(map[string]*tview.TextView)
	chanPages := tview.NewPages()

	for i := 0; i < listNum; i++ {//initalize pages on per-channel basis
		name := fmt.Sprintf("channel %d", i)
		view, prim := channel(name, onSend)//channel returns TextView and Primitive
		msgView[name] = view
			chanPages.AddPage(name,prim,true,i == 0)
		}
	
	list := tview.NewList()//initalize list
	for i := 0; i < listNum; i++ {
		name := fmt.Sprintf("channel %d", i)
		list.AddItem(name,"",0, func(){
		activeChannel = name
		chanPages.SwitchToPage(name)//switches chan pages, not page
		})
	}


chat := tview.NewFlex().SetDirection(tview.FlexRow).//adds everything to flex, to be refrenced later
	AddItem(chanPages,0,1, true)
mainView := tview.NewFlex().
	AddItem(list,20,0,false).
	AddItem(chat,0,1,true)

	signIn := tview.NewForm()//input form, ran on launch
		signIn.AddInputField("username", "", 20, nil, nil)
		signIn.AddInputField("password", "", 20, nil, nil)
		signIn.AddDropDown("username color",[]string{"red","blue","green"},0,nil)
		signIn.AddButton("enter", func() {
		currentUser = signIn.GetFormItemByLabel("username").(*tview.InputField).GetText()
		currentPass = signIn.GetFormItemByLabel("password").(*tview.InputField).GetText()
		_, color := signIn.GetFormItemByLabel("username color").(*tview.DropDown).GetCurrentOption()
		_ = color //negate unused var
			//fmt.Println(currentUser,currentPass,color)
			go listenClient()
			//go runClientListener()
			pages.SwitchToPage("main")//switch to other pages
			app.SetFocus(list)//focus to whatever is selected
		})
	signIn.SetBorder(true).SetTitle("enter details").SetTitleAlign(tview.AlignLeft)
	
	pages.AddPage("signIn",signIn,true,true)
	pages.AddPage("main",mainView,true,false)//unable to be reached unless with sign in form

app.SetRoot(pages,true).SetFocus(signIn) //puts the user in sign in to start
}

func channel(name string, onSend func(string)) (*tview.TextView, tview.Primitive){
	textView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)
		textView.ScrollToEnd()
	textView.SetBackgroundColor(tcell.ColorDefault)
	textView.SetBorderPadding(0, 0, 1, 1)
	textView.SetBorder(true)
	textView.SetTitle(" PingChat v2 ")

		inputBox := tview.NewInputField()
	inputBox.SetBorder(true)
	inputBox.SetBackgroundColor(tcell.ColorDefault)
	inputBox.SetFieldBackgroundColor(tcell.ColorDefault)

	inputBox.SetDoneFunc(func(key tcell.Key) {//take input from inputBox
		if key != tcell.KeyEnter {
			return
		}
		text := inputBox.GetText()
		if text == "" || len(text) > 512 {
			return
		}
		inputBox.SetText("")
		go onSend(text)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).//puts text view and input into respective positions
		AddItem(textView, 0, 1, false).
		AddItem(inputBox, 3, 1, true)

	return textView, flex
}

func tuiPrint(line string) {
	view := msgView[activeChannel]//updates with new messages
	app.QueueUpdateDraw(func() {
		fmt.Fprintf(view, "%s\n", line)
	})
}

func runTUI() {//throwing errors
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func ScrollToMessage(){
	textView := msgView[activeChannel]
	app.QueueUpdateDraw(func(){
		textView.ScrollToEnd()
	})
}
