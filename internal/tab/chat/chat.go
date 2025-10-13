package chat

import (
	"context"
	"fmt"
	"io"
	"log"
	"vado-client/internal/appcontext"
	pb "vado-client/internal/pb/chat"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewChat(appCtx *appcontext.AppContext) *fyne.Container {
	client := pb.NewChatServiceClient(appCtx.GRPC)
	ctx, _ := context.WithCancel(context.Background())

	//win.SetCloseIntercept(func() {
	//	fmt.Println("CLOSE!")
	//cancel()
	//})

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
			dialog.ShowError(err, appCtx.Win)
			return
		}
		input.SetText("")
	})

	content := container.NewBorder(
		container.NewVBox(nameEntry, messages),
		container.NewVBox(input, sendButton),
		nil, nil,
	)

	// поток сообщений
	go func() {
		stream, err := client.ChatStream(ctx, &pb.Empty{})
		if err != nil {
			dialog.ShowError(err, appCtx.Win)
			return
		}

		for {
			msg, err := stream.Recv()
			if err == io.EOF || ctx.Err() != nil {
				break
			}
			if err != nil {
				log.Printf("Ошибка получения: %v", err)
				break
			}
			text := fmt.Sprintf("%s: %s\n", msg.User, msg.Text)
			fyne.Do(func() {
				messages.SetText(messages.Text + text)
			})
		}
	}()

	return content
}
