package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gcerrato/godog/src/llm"
)

type MessageType string

const (
	MessageTypeUser      MessageType = "user"
	MessageTypeSystem    MessageType = "system"
	MessageTypeAssistant MessageType = "assistant"
)

type Message struct {
	Text string
	MessageType
}

func GetChatBubble(message Message) *fyne.Container {
	label := widget.NewLabel(fmt.Sprintf("%s: %s", message.MessageType, message.Text))
	//label.Wrapping = fyne.TextWrapWord

	copyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		clipboard.SetContent(message.Text)
	})

	if message.MessageType == MessageTypeSystem || message.MessageType == MessageTypeAssistant {
		return container.NewHBox(label, copyButton, layout.NewSpacer())
	} else {
		return container.NewHBox(layout.NewSpacer(), label, copyButton)
	}
}

func AddMessage(ctx context.Context, message Message, chatContainer *fyne.Container, textEntry *widget.Entry) {
	textEntry.SetText("")
	loadingLabel := widget.NewLabel("Assistant is thinking...")

	if message.Text != "" {
		chatContainer.Add(GetChatBubble(message))
	}

	chatContainer.Add(loadingLabel)
	response, err := llm.SendToLLM(ctx, message.Text)
	if err != nil {
		chatContainer.Remove(loadingLabel)
		chatContainer.Add(GetChatBubble(Message{Text: "Error", MessageType: MessageTypeSystem}))
		return
	}
	chatContainer.Remove(loadingLabel)
	chatContainer.Add(GetChatBubble(Message{Text: response, MessageType: MessageTypeAssistant}))

}

func main() {
	a := app.New()
	w := a.NewWindow("Chat Window")
	ctx := context.Background()

	chatContainer := container.NewVBox()
	w.Resize(fyne.NewSize(500, 500))

	chatContainer.Add(GetChatBubble(Message{Text: "Hi!", MessageType: MessageTypeSystem}))

	textEntry := widget.NewEntry()
	textEntry.OnSubmitted = func(text string) {
		if text != "" {
			AddMessage(ctx, Message{Text: textEntry.Text, MessageType: MessageTypeUser}, chatContainer, textEntry)
		}
	}

	loadButton := widget.NewButton("Load", func() {

	})
	sendButton := widget.NewButton("Send", func() {
		message := Message{Text: textEntry.Text, MessageType: MessageTypeUser}
		AddMessage(ctx, message, chatContainer, textEntry)
	})

	inputContainer := container.NewBorder(nil, nil, loadButton, sendButton, textEntry)
	scrollContainer := container.NewScroll(chatContainer)
	w.SetContent(container.NewBorder(nil, inputContainer, nil, nil, scrollContainer))
	w.ShowAndRun()
}
