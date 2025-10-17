package chat

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MessageItem struct {
	widget.BaseWidget

	label *widget.Label
}

func NewMessageItem() *MessageItem {
	item := &MessageItem{
		label: widget.NewLabel(""),
	}

	item.ExtendBaseWidget(item)

	return item
}

func (item *MessageItem) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewHBox(item.label)
	return widget.NewSimpleRenderer(content)
}

func (item *MessageItem) SetData(text string) {
	item.label.SetText(text)
}
