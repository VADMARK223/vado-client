package chat

import (
	"log"
	pb "vado-client/internal/pb/chat"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc"
)

func NewChat() {
	// подключаемся к gRPC серверу
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("не удалось подключиться: %v", err)
	}

	client := pb.NewChatServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())

	// При закрытии окна — корректно завершаем соединение
	w.SetCloseIntercept(func() {
		cancel()
		conn.Close()
		a.Quit()
	})

	messages := widget.NewMultiLineEntry()
	messages.Disable()

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Ваше имя")

	input := widget.NewEntry()
	input.SetPlaceHolder("Сообщение...")

	sendButton := widget.NewButton("Отправить", func() {
		user := nameEntry.Text
		text := input.Text
		if user == "" || text == "" {
			return
		}

		_, err := client.SendMessage(ctx, &pb.ChatMessage{
			User: user,
			Text: text,
		})
		if err != nil {
			log.Printf("Ошибка отправки: %v", err)
			return
		}
		input.SetText("")
	})

	content := container.NewBorder(
		container.NewVBox(nameEntry, messages),
		container.NewVBox(input, sendButton),
		nil, nil,
	)
}
