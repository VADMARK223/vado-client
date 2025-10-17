package chat

import (
	"context"
	"fmt"
	"io"
	"log"
	pb "vado-client/api/pb/chat"
	"vado-client/internal/appcontext"
	"vado-client/internal/component/userInfo"
	"vado-client/internal/grpc/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc/metadata"
)

func withAuth(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

func NewChat(appCtx *appcontext.AppContext, a fyne.App) *fyne.Container {
	clientGRPC := pb.NewChatServiceClient(appCtx.GRPC)
	ctx, _ := context.WithCancel(context.Background())

	//win.SetCloseIntercept(func() {
	//	fmt.Println("CLOSE!")
	//cancel()
	//})

	messages := widget.NewMultiLineEntry()
	messages.Disable()

	input := widget.NewEntry()
	input.SetPlaceHolder("Сообщение...")

	sendButton := widget.NewButton("Отправить", func() {
		text := input.Text
		token := client.GetToken(a)
		appCtx.Log.Debugf("Send with token: %s", token)
		authCtx := withAuth(ctx, token)
		_, err := clientGRPC.SendMessage(authCtx, &pb.ChatMessage{
			User: client.GetUsername(a),
			Text: text,
		})
		if err != nil {
			dialog.ShowError(err, appCtx.Win)
			return
		}
		input.SetText("")
	})

	loginBtn := widget.NewButton("Вход", func() {
		userInfo.ShowLoginDialog(appCtx, a, nil)
	})

	updateButtons(a, sendButton, loginBtn)

	userNameText := widget.NewRichTextFromMarkdown(fmt.Sprintf("Привет, **%s**!", client.GetUsername(a)))
	a.Preferences().AddChangeListener(func() {
		userNameText.ParseMarkdown(fmt.Sprintf("Привет, **%s**!", client.GetUsername(a)))
		userNameText.Refresh()
		updateButtons(a, sendButton, loginBtn)
	})

	content := container.NewBorder(
		container.NewVBox(userNameText, messages),
		container.NewVBox(input, sendButton, loginBtn),
		nil, nil,
	)

	// поток сообщений
	go func() {
		stream, err := clientGRPC.ChatStream(ctx, &pb.Empty{})
		if err != nil {
			appCtx.Log.Errorw("Ошибка создания потока", "error", err.Error())
			//dialog.ShowInformation("Ошибка создания потока", err.Error(), appCtx.Win)
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

func updateButtons(a fyne.App, sendBtn *widget.Button, loginBtn *widget.Button) {
	if client.IsAuth(a) {
		sendBtn.Show()
		loginBtn.Hide()
	} else {
		sendBtn.Hide()
		loginBtn.Show()
	}
}
