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

	controlBox := container.NewHBox(widget.NewButton("1", func() {}), widget.NewButton("2", func() {}))
	controlBox.Hide()
	loginBtn := widget.NewButton("Вход", func() {
		userInfo.ShowLoginDialog(appCtx, nil)
	})

	userNameText := widget.NewRichTextFromMarkdown(fmt.Sprintf("Привет, **%s**!", client.GetUsername(appCtx.App)))

	topBox := container.NewVBox(controlBox, userNameText)

	messages := binding.NewUntypedList()
	/*_ = messages.Set([]interface{}{
		"apple",
		"banana",
		"cherry",
		"cherry1",
		"cherry2",
		"cherry43",
		"cherry5",
		"cherry6",
		"cherry7",
		"cherry8",
		"cherry9",
	})*/
	list := widget.NewListWithData(
		messages,
		func() fyne.CanvasObject { return NewMessageItem() },
		func(item binding.DataItem, obj fyne.CanvasObject) {
			str, _ := item.(binding.Untyped).Get()
			//obj.(*widget.Label).SetText(str.(string))

			messageItem := obj.(*MessageItem)
			messageItem.SetData(str.(string))
		})

	scroll := container.NewVScroll(list)

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
		_, errSendMessage := clientGRPC.SendMessage(authCtx, &pb.ChatMessage{
			User: client.GetUsername(appCtx.App),
			Text: text,
		})
		if errSendMessage != nil {
			dialog.ShowError(errSendMessage, appCtx.Win)
			return
		}
		input.SetText("")
	})

	updateButtons(appCtx.App, sendButton, loginBtn)

	bottomBox := container.NewVBox(input, sendButton)

	content := container.NewBorder(topBox, bottomBox, nil, nil, scroll)

	// Поток сообщений
	go func() {
		defer cancel()

		for {
			if ctx.Err() != nil {
				return
			}

			token := client.GetToken(appCtx.App)
			authCtx := withAuth(ctx, token)

			stream, errStream := clientGRPC.ChatStream(authCtx, &pb.Empty{})
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
					errAppend := messages.Append(msg.Text)
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
		updateButtons(appCtx.App, sendButton, loginBtn)
	})

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
