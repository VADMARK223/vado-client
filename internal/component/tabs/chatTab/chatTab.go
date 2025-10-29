package chatTab

import (
	"fmt"
	"time"
	pb "vado-client/api/pb/chat"
	"vado-client/internal/app"
	"vado-client/internal/app/keyman"
	"vado-client/internal/component/common/userInfo"
	"vado-client/internal/component/tabs/tabItem"

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

func New(appCtx *app.Context) tabItem.TabContent {
	clientGRPC := pb.NewChatServiceClient(appCtx.GRPC)

	infoText := widget.NewRichTextFromMarkdown("Необходимо выполнить вход.")
	userCountTxt := widget.NewRichTextWithText("")
	updateCountText(userCountTxt, 0)
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

	scrollToDown := func() {
		fyne.Do(func() {
			list.Refresh()
			n := messages.Length()
			if n > 0 {
				list.ScrollTo(n - 1)
			}
		})
	}

	scrollDownBtn := widget.NewButton("Вниз", func() {
		scrollToDown()
	})
	scrollDownBtn.Resize(fyne.NewSize(45, 40))

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			fyne.Do(func() {
				scrollDownBtn.Move(fyne.NewPos(
					list.Size().Width-scrollDownBtn.Size().Width-8,
					list.Size().Height-scrollDownBtn.Size().Height-8))
				totalItems := list.Length()
				if totalItems == 0 {
					scrollDownBtn.Hide()
					return
				}

				itemHeight := list.MinSize().Height
				itemsH := itemHeight * float32(totalItems)

				if itemsH > list.GetScrollOffset()+list.Size().Height {
					scrollDownBtn.Show()
				} else {
					scrollDownBtn.Hide()
				}
			})
		}
	}()

	messages.AddListener(binding.NewDataListener(func() {
		scrollToDown()
	}))

	loginBtn := widget.NewButton("Вход", func() {
		userInfo.ShowLoginDialog(appCtx, nil)
	})

	updateButtons(appCtx, loginBtn)

	input, sendBtn := newInputBox(appCtx, clientGRPC)
	inputBox := container.NewVBox(input, sendBtn)

	scrollWithBtn := container.NewStack(list, container.NewWithoutLayout(scrollDownBtn))
	scrollDownBtn.Hide()

	content := container.NewBorder(container.NewVBox(infoText, userCountTxt), inputBox, nil, nil, scrollWithBtn)

	manager := NewChatStreamManager(appCtx, clientGRPC, messages, func(count uint32) {
		updateCountText(userCountTxt, count)
	})

	updateVisibility := func() {
		isAuth := appCtx.Prefs.IsAuth()
		appCtx.Log.Infow("Change listeners", "isAuth", isAuth)
		if isAuth {
			infoText.Hide()
			userCountTxt.Show()
			loginBtn.Hide()
			inputBox.Show()
		} else {
			infoText.Show()
			userCountTxt.Hide()
			loginBtn.Show()
			inputBox.Hide()
		}
	}

	appCtx.Prefs.ChangeListeners(func() {
		fyne.Do(func() {
			updateVisibility()
		})

		if appCtx.Prefs.IsAuth() {
			manager.Start()
		} else {
			manager.Stop()
		}
	})
	updateVisibility()

	if appCtx.Prefs.IsAuth() {
		manager.Start()
	}

	return &ChatTab{
		Container: content,
		btn:       sendBtn,
		canvas:    appCtx.Win.Canvas(),
		keyman:    appCtx.KeyMan,
	}
}

func updateCountText(txt *widget.RichText, count uint32) {
	if count == 1 {
		txt.ParseMarkdown("В комнате пока только вы.")
	} else {
		txt.ParseMarkdown(fmt.Sprintf("В комнате пользователей: **%v**", count))
	}
}

func updateButtons(appCtx *app.Context, loginBtn *widget.Button) {
	if appCtx.Prefs.IsAuth() {
		loginBtn.Hide()
	} else {
		loginBtn.Show()
	}
}
