package kafkaTab

import (
	"context"
	"fmt"
	"log"
	"time"
	"vado-client/internal/app"
	"vado-client/internal/component/tabs/tabItem"
	"vado-client/internal/grpc/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/segmentio/kafka-go"
)

type Tab struct {
	*fyne.Container
	cxt *app.Context
}

func (t *Tab) Open() {
	t.cxt.Log.Debugw("Kafka tab opened")
}

func (t *Tab) Close() {
	t.cxt.Log.Debugw("Kafka tab closed")
}

func (t *Tab) Canvas() fyne.CanvasObject {
	return t.Container
}

func New(appCtx *app.Context) tabItem.TabContent {
	appCtx.Log.Debugw("Kafka tab created")
	writer := newKafkaWriter("localhost:9094", "chat")

	input := widget.NewEntry()
	input.SetPlaceHolder("Введите сообщение")
	time.AfterFunc(100*time.Millisecond, func() {
		fyne.Do(func() {
			fyne.CurrentApp().Driver().CanvasForObject(input).Focus(input)
		})
	})

	btn := widget.NewButton("Отправить", func() {
		appCtx.Log.Debugw("Kafka tab send message")
		user := fmt.Sprintf("User-%s", client.GetUsername(appCtx.App))
		msg := input.Text
		if msg == "" {
			dialog.ShowInformation("Предупреждение", "Пустое сообщение", appCtx.Win)
			return
		}
		sendChatMessage(writer, user, msg)
	})
	btn.Disable()

	input.OnSubmitted = func(text string) {
		btn.OnTapped()
	}
	input.OnChanged = func(text string) {
		if text != "" {
			btn.Enable()
		} else {
			btn.Disable()
		}
	}

	c := container.NewVBox(
		input,
		btn,
	)

	return &Tab{
		Container: c,
		cxt:       appCtx,
	}
}

func newKafkaWriter(broker, topic string) *kafka.Writer {
	brokers := []string{broker}
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{}, // Балансировщик для распределения сообщений по партициям (можно использовать другие: Hash, RoundRobin)
		AllowAutoTopicCreation: true,                // Авто создание топика
	}

	return writer
}

func sendChatMessage(writer *kafka.Writer, user, msg string) {
	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(user),
			Value: []byte(msg),
		},
	)
	if err != nil {
		log.Printf("Kafka write error: %v", err)
	}
}
