package chat

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"
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

func NewChat(appCtx *appcontext.AppContext) *fyne.Container {
	clientGRPC := pb.NewChatServiceClient(appCtx.GRPC)
	ctx, cancel := context.WithCancel(context.Background())

	messages := widget.NewMultiLineEntry()
	messages.Disable()

	input := widget.NewEntry()
	input.SetPlaceHolder("Сообщение...")

	sendButton := widget.NewButton("Отправить", func() {
		text := input.Text
		if text == "" {
			dialog.ShowInformation("Предупреждение", "Пустое сообщение", appCtx.Win)
			return
		}
		token := client.GetToken(appCtx.App)
		appCtx.Log.Debugf("Send with token: %s", token)
		authCtx := withAuth(ctx, token)
		_, err := clientGRPC.SendMessage(authCtx, &pb.ChatMessage{
			User: client.GetUsername(appCtx.App),
			Text: text,
		})
		if err != nil {
			dialog.ShowError(err, appCtx.Win)
			return
		}
		input.SetText("")
	})

	loginBtn := widget.NewButton("Вход", func() {
		userInfo.ShowLoginDialog(appCtx, nil)
	})

	updateButtons(appCtx.App, sendButton, loginBtn)

	userNameText := widget.NewRichTextFromMarkdown(fmt.Sprintf("Привет, **%s**!", client.GetUsername(appCtx.App)))
	appCtx.App.Preferences().AddChangeListener(func() {
		userNameText.ParseMarkdown(fmt.Sprintf("Привет, **%s**!", client.GetUsername(appCtx.App)))
		userNameText.Refresh()
		updateButtons(appCtx.App, sendButton, loginBtn)
	})

	content := container.NewBorder(
		container.NewVBox(userNameText, messages),
		container.NewVBox(input, sendButton, loginBtn),
		nil, nil,
	)

	// Поток сообщений
	go func() {
		defer cancel() // безопасность — отмена при выходе из горутины
		var builder strings.Builder

		for {
			fmt.Println("TICK")
			if ctx.Err() != nil {
				return
			}

			token := client.GetToken(appCtx.App)
			authCtx := withAuth(ctx, token)

			stream, err := clientGRPC.ChatStream(authCtx, &pb.Empty{})
			if err != nil {
				appCtx.Log.Errorw("Ошибка подключения к потоку", "error", err)
				time.Sleep(2 * time.Second)
				continue
			}

			appCtx.Log.Infow("Подключен к потоку сообщений")

			for {
				msg, err := stream.Recv()
				if err == io.EOF || ctx.Err() != nil {
					appCtx.Log.Infow("Завершение потока")
					break
				}
				if err != nil {
					appCtx.Log.Errorw("Ошибка получения сообщения", "error", err)
					break
				}

				builder.WriteString(fmt.Sprintf("%s: %s\n", msg.User, msg.Text))

				fyne.Do(func() {
					messages.SetText(builder.String())
				})
			}

			// Если был обрыв — попробуем переподключиться через 2 секунды
			time.Sleep(2 * time.Second)
		}
	}()

	appCtx.AddCloseHandler(func() {
		cancel()
	})

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
