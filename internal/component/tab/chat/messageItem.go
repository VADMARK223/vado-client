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
	usernameLbl *widget.Label
	timeLbl     *widget.Label
	messageLbl  *widget.Label
	isMyMessage bool
}

func NewMessageItem() *MessageItem {
	item := &MessageItem{
		usernameLbl: widget.NewLabel(""),
		timeLbl:     widget.NewLabel(""),
		messageLbl:  widget.NewLabel(""),
	}

	item.usernameLbl.TextStyle = fyne.TextStyle{Bold: true}

	item.messageLbl.TextStyle = fyne.TextStyle{Italic: true}
	item.messageLbl.Wrapping = fyne.TextWrapWord

	item.timeLbl.TextStyle = fyne.TextStyle{Monospace: true}

	item.ExtendBaseWidget(item)
	return item
}

func (item *MessageItem) CreateRenderer() fyne.WidgetRenderer {
	// Заголовок: имя пользователя и время в одной строке
	header := container.NewHBox(
		item.usernameLbl,
		layout.NewSpacer(),
		item.timeLbl,
	)

	// Контент: заголовок и текст сообщения
	content := container.NewVBox(
		header,
		item.messageLbl,
	)

	// Добавляем отступы вокруг всего сообщения
	paddedContent := container.NewPadded(content)

	return widget.NewSimpleRenderer(paddedContent)
}

func (item *MessageItem) SetData(data *chat.ChatMessage, isMyMessage bool) {
	item.usernameLbl.SetText(data.GetUser())
	//item.timeLbl.SetText(formatTime(data.GetCreatedAt()))
	item.timeLbl.SetText("time")
	item.messageLbl.SetText(data.GetText())
	item.isMyMessage = isMyMessage

	if isMyMessage {
		item.usernameLbl.Importance = widget.LowImportance
		item.usernameLbl.Alignment = fyne.TextAlignTrailing
		item.timeLbl.Alignment = fyne.TextAlignTrailing
		item.messageLbl.Alignment = fyne.TextAlignTrailing
	} else {
		item.usernameLbl.Importance = widget.HighImportance
		item.usernameLbl.Alignment = fyne.TextAlignLeading
		item.timeLbl.Alignment = fyne.TextAlignLeading
		item.messageLbl.Alignment = fyne.TextAlignLeading
	}

	item.Refresh()
}

func formatTime(timestamp string) string {
	// Реализуйте форматирование времени
	return timestamp
}
