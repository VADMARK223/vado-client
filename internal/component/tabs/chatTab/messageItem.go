package chatTab

import (
	"fmt"
	"image/color"
	"time"
	"vado-client/api/pb/chat"
	"vado-client/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MessageItem struct {
	widget.BaseWidget
	usernameLbl *canvas.Text
	timeLbl     *widget.Label
	messageLbl  *widget.Label
	isMyMessage bool
}

func NewMessageItem() *MessageItem {
	item := &MessageItem{
		usernameLbl: canvas.NewText("", color.White),
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

func (item *MessageItem) SetData(data *chat.ChatMessage) {
	item.usernameLbl.Text = data.User.Username
	item.usernameLbl.Color = hexToColor(data.User.Color)
	isMyMessage := data.Type == chat.MessageType_MESSAGE_SELF
	item.timeLbl.SetText(formatTime(data.Timestamp))
	item.messageLbl.SetText(data.GetText())
	item.isMyMessage = isMyMessage

	if isMyMessage {
		//item.usernameLbl.Importance = widget.LowImportance
		item.usernameLbl.Alignment = fyne.TextAlignTrailing
		item.timeLbl.Alignment = fyne.TextAlignTrailing
		item.messageLbl.Alignment = fyne.TextAlignTrailing
	} else {
		//item.usernameLbl.Importance = widget.HighImportance
		item.usernameLbl.Alignment = fyne.TextAlignLeading
		item.timeLbl.Alignment = fyne.TextAlignLeading
		item.messageLbl.Alignment = fyne.TextAlignLeading
	}

	item.Refresh()
}

func formatTime(ts int64) string {
	t := time.Unix(ts, 0).Local()
	return utils.FormatTime(t)
}

func hexToColor(hex string) color.Color {
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return theme.Color(theme.ColorNameForeground)
	}
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}
