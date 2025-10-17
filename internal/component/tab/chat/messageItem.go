package chat

import (
	"vado-client/api/pb/chat"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MessageItem struct {
	widget.BaseWidget

	username *widget.Label
	message  *widget.Label
	content  *fyne.Container
	spacer   fyne.CanvasObject
}

func NewMessageItem() *MessageItem {
	item := &MessageItem{
		username: widget.NewLabel(""),
		message:  widget.NewLabel(""),
		spacer:   layout.NewSpacer(),
	}

	item.ExtendBaseWidget(item)

	return item
}

func (item *MessageItem) CreateRenderer() fyne.WidgetRenderer {
	item.content = container.NewHBox(item.spacer, item.username, item.message)
	return widget.NewSimpleRenderer(item.content)
}

func (item *MessageItem) SetData(data *chat.ChatMessage, isMyMessage bool) {
	item.username.SetText(data.GetUser())
	item.message.SetText(data.GetText())

	if isMyMessage {
		item.username.Hide()
		item.spacer.Hide()
		item.message.TextStyle = fyne.TextStyle{Bold: true}
	} else {
		item.username.TextStyle = fyne.TextStyle{Italic: true}
		item.message.TextStyle = fyne.TextStyle{Italic: true}
	}

	item.Refresh()
}
