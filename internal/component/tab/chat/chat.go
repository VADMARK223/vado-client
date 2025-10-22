package chat

import (
	"context"
	"fmt"
	"io"
	"time"
	pb "vado-client/api/pb/chat"
	"vado-client/internal/appcontext"
	"vado-client/internal/component/userInfo"
	"vado-client/internal/grpc/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc/metadata"
)

func withAuth(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

func NewChat(appCtx *appcontext.AppContext) *fyne.Container {
	clientGRPC := pb.NewChatServiceClient(appCtx.GRPC)
	ctx, cancel := context.WithCancel(context.Background())

	controlBox := container.NewHBox(widget.NewButton("1", func() {}), widget.NewButton("2", func() {}))
	controlBox.Hide()
	loginBtn := widget.NewButton("Вход", func() {
		userInfo.ShowLoginDialog(appCtx, nil)
	})

	userNameText := widget.NewRichTextFromMarkdown(fmt.Sprintf("Привет, **%s**!", client.GetUsername(appCtx.App)))

	topBox := container.NewVBox(controlBox, userNameText)

	messages := binding.NewUntypedList()
	list := widget.NewListWithData(
		messages,
		func() fyne.CanvasObject { return NewMessageItem() },
		func(item binding.DataItem, obj fyne.CanvasObject) {
			str, _ := item.(binding.Untyped).Get()
			messageItem := obj.(*MessageItem)
			messageData := str.(*pb.ChatMessage)
			messageItem.SetData(messageData)
		})

	scroll := container.NewVScroll(list)

	updateButtons(appCtx.App, loginBtn)

	content := container.NewBorder(topBox, newInputBox(appCtx, ctx, clientGRPC), nil, nil, scroll)

	// Поток сообщений
	go func() {
		defer cancel()

		for {
			if ctx.Err() != nil {
				return
			}

			token := client.GetToken(appCtx.App)
			authCtx := withAuth(ctx, token)

			userID := client.GetUserID(appCtx.App)

			req := &pb.ChatStreamRequest{Id: userID}

			stream, errStream := clientGRPC.ChatStream(authCtx, req)
			if errStream != nil {
				appCtx.Log.Errorw("Ошибка подключения к потоку", "error", errStream)
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

				fyne.Do(func() {
					errAppend := messages.Append(msg)
					if errAppend != nil {
						appCtx.Log.Errorw("Error append message", "error", errAppend)
					}
				})
			}

			// Если был обрыв — попробуем переподключиться через 2 секунды
			time.Sleep(2 * time.Second)
		}
	}()

	appCtx.App.Preferences().AddChangeListener(func() {
		userNameText.ParseMarkdown(fmt.Sprintf("Привет, **%s**!", client.GetUsername(appCtx.App)))
		userNameText.Refresh()
		updateButtons(appCtx.App, loginBtn)
	})

	appCtx.AddCloseHandler(func() {
		cancel()
	})

	return content
}

func updateButtons(a fyne.App, loginBtn *widget.Button) {
	if client.IsAuth(a) {
		loginBtn.Hide()
	} else {
		loginBtn.Show()
	}
}
