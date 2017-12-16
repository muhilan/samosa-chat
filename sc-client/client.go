package main

import (
    "github.com/andlabs/ui"
)

func main() {
    err := ui.Main(func() {
        chatTxt := ui.NewEntry()
        button := ui.NewButton("Send")
        greeting := ui.NewLabel("")
        box := ui.NewVerticalBox()
        // box.Append(ui.NewLabel("Enter your name:"), false)
        box.Append(chatTxt, false)
        box.Append(button, false)
        box.Append(greeting, false)
        window := ui.NewWindow("Samosa-chat", 200, 100, false)
        window.SetChild(box)
        button.OnClicked(func(*ui.Button) {
            greeting.SetText("\n" + chatTxt.Text() )
        })
        window.OnClosing(func(*ui.Window) bool {
            ui.Quit()
            return true
        })
        window.Show()
    })
    if err != nil {
        panic(err)
    }
}