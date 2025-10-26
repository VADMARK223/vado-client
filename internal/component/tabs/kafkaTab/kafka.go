package kafkaTab

import (
	"context"
	"fmt"
	"vado-client/internal/app"
	"vado-client/internal/component/tabs/tabItem"
	"vado-client/internal/grpc/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Tab struct {
	*fyne.Container
	cxt      *app.Context
	input    *widget.Entry
	btn      *widget.Button
	producer *Producer
}

func (t *Tab) Open() {
	t.cxt.Log.Debugw("Kafka tab opened")
	producer := NewProducer("localhost:9094", "chat", t.cxt.Log)
	t.producer = producer

	t.btn.Enable()
	t.btn.OnTapped = func() {
		t.sendMessage()
	}
}

func (t *Tab) Close() {
	t.cxt.Log.Debugw("Kafka tab closed")
	if t.producer != nil {
		if err := t.producer.Close(); err != nil {
			t.cxt.Log.Errorw("Kafka producer close error", "error", err)
		}
		t.producer = nil
	}
	t.btn.Disable()
}

func (t *Tab) Canvas() fyne.CanvasObject {
	return t.Container
}

func New(appCtx *app.Context) tabItem.TabContent {
	appCtx.Log.Debugw("Kafka tab created")

	input := widget.NewEntry()
	input.SetPlaceHolder("Введите сообщение")

	btn := widget.NewButton("Отправить", nil)
	btn.Disable()

	tab := &Tab{
		cxt:   appCtx,
		input: input,
		btn:   btn,
	}

	input.OnChanged = func(text string) {
		if text != "" {
			btn.Enable()
		} else {
			btn.Disable()
		}
	}
	input.OnSubmitted = func(string) {
		tab.sendMessage()
	}

	tab.Container = container.NewVBox(input, btn)
	return tab
}

// sendMessage — логика отправки одного сообщения
func (t *Tab) sendMessage() {
	user := fmt.Sprintf("User-%s", client.GetUsername(t.cxt.App))
	msg := t.input.Text

	if msg == "" {
		dialog.ShowInformation("Предупреждение", "Пустое сообщение", t.cxt.Win)
		return
	}
	if t.producer == nil {
		dialog.ShowError(fmt.Errorf("producer ещё не инициализирован"), t.cxt.Win)
		return
	}

	err := t.producer.SendMessage(context.Background(), []byte(user), []byte(msg))
	if err != nil {
		dialog.ShowError(fmt.Errorf("не удалось отправить сообщение: %v", err), t.cxt.Win)
		return
	}

	t.input.SetText("")
}
