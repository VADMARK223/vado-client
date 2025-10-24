package chatTab

import (
	"context"
	"fmt"
	"io"
	"time"
	pb "vado-client/api/pb/chat"
	"vado-client/internal/app"
	"vado-client/internal/app/keyman"
	"vado-client/internal/component/tabs/tabItem"
	"vado-client/internal/component/userInfo"
	"vado-client/internal/grpc/client"
	"vado-client/internal/grpc/middleware"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ChatTab struct {
	*fyne.Container
	btn    *widget.Button
	canvas fyne.Canvas
	unsub  func() // Отписка от события
	keyman *keyman.KeyManager
}

func (t *ChatTab) Open() {
	t.unsub = t.keyman.Subscribe(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
			t.btn.OnTapped()
		}
	})
}

func (t *ChatTab) Close() {
	if t.unsub != nil {
		t.unsub()
	}
}

func (t *ChatTab) Canvas() fyne.CanvasObject {
	return t.Container
}

var userCountText = widget.NewRichTextWithText("")

func New(appCtx *app.Context) tabItem.TabContent {
	clientGRPC := pb.NewChatServiceClient(appCtx.GRPC)
	ctx, cancel := context.WithCancel(context.Background())

	controlBox := container.NewHBox(widget.NewButton("1", func() {}), widget.NewButton("2", func() {}))
	controlBox.Hide()
	loginBtn := widget.NewButton("Вход", func() {
		userInfo.ShowLoginDialog(appCtx, nil)
	})

	updateCountText(0)

	topBox := container.NewVBox(controlBox, userCountText)

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

	input, sendBtn := newInputBox(appCtx, ctx, clientGRPC)
	inputBox := container.NewVBox(input, sendBtn)

	content := container.NewBorder(topBox, inputBox, nil, nil, scroll)

	// Поток сообщений
	go func() {
		defer cancel()

		for {
			if ctx.Err() != nil {
				return
			}

			userID := client.GetUserID(appCtx.App)

			req := &pb.ChatStreamRequest{Id: userID}
			stream, errStream := clientGRPC.ChatStream(middleware.WithAuth(appCtx, ctx), req)
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

					//_, _ = pp.Println(msg)
					if msg.Type == pb.MessageType_MESSAGE_SYSTEM {
						updateCountText(msg.UsersCount)
					}
				})
			}

			// Если был обрыв — попробуем переподключиться через 2 секунды
			time.Sleep(2 * time.Second)
		}
	}()

	appCtx.App.Preferences().AddChangeListener(func() {
		updateButtons(appCtx.App, loginBtn)
	})

	appCtx.AddCloseHandler(func() {
		cancel()
	})

	return &ChatTab{
		Container: content,
		btn:       sendBtn,
		canvas:    appCtx.Win.Canvas(),
		keyman:    appCtx.KeyMan,
	}
}

func updateCountText(count uint32) {
	if count == 1 {
		userCountText.ParseMarkdown("В комнате пока только вы.")
	} else {
		userCountText.ParseMarkdown(fmt.Sprintf("В комнате пользователей: **%v**", count))
	}
}

func updateButtons(a fyne.App, loginBtn *widget.Button) {
	if client.IsAuth(a) {
		loginBtn.Hide()
	} else {
		loginBtn.Show()
	}
}
