package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gcerrato/godog/src/llm"
)


type Message struct {
	Text   string
	IsSystem bool
}

func GetChatBubble(message Message) *fyne.Container {
	if message.IsSystem {
		return container.NewHBox(widget.NewLabel(fmt.Sprintf("system: %s", message.Text)), layout.NewSpacer())
	} else {

		return container.NewHBox(layout.NewSpacer(), widget.NewLabel(fmt.Sprintf("user: %s", message.Text)))
	}
}

func AddMessage(ctx context.Context, message Message, chatContainer *fyne.Container, textEntry *widget.Entry) {
		if message.Text != "" {
			chatContainer.Add(GetChatBubble(message))
		}
		response, err := llm.SendToLLM(ctx,message.Text)
		if err != nil {
			chatContainer.Add(GetChatBubble(Message{Text: "Error", IsSystem: true}))

		}
		chatContainer.Add(GetChatBubble(Message{Text: response, IsSystem: true}))

	textEntry.SetText("")
}

func main() {
	a := app.New()
	w := a.NewWindow("Chat Window")
	ctx := context.Background()


	chatContainer := container.NewVBox()
	w.Resize(fyne.NewSize(500, 500))

	chatContainer.Add(GetChatBubble(Message{Text: "Hi!", IsSystem: true}))

	textEntry := widget.NewEntry()
	textEntry.OnSubmitted = func(text string) {
		if text != "" {
			AddMessage(ctx, Message{Text: textEntry.Text, IsSystem: false}, chatContainer, textEntry)
		}
	}

	loadButton := widget.NewButton("Load", func() {

	})
	sendButton := widget.NewButton("Send", func() {
		message := Message{Text: textEntry.Text, IsSystem: false}
		AddMessage(ctx, message, chatContainer, textEntry)
	})

	inputContainer := container.NewBorder(nil, nil, loadButton, sendButton, textEntry)
	w.SetContent(container.NewBorder(nil, inputContainer, nil, nil, chatContainer))
	w.ShowAndRun()
}
