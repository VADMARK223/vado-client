package chatTab

import (
	"context"
	"time"
	pb "vado-client/api/pb/chat"
	"vado-client/internal/app"
	"vado-client/internal/component/common"
	"vado-client/internal/grpc/middleware"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func newInputBox(ctx *app.Context, clientGRPC pb.ChatServiceClient, sendBtn *widget.Button) *fyne.Container {
	msgInput := common.NewVadoEntry()
	msgInput.SetMaxLen(30)
	msgInput.SetText(ctx.Prefs.LastInput())
	msgInput.SetPlaceHolder("White a message...")

	sendBtn.OnTapped = func() {
		text := msgInput.Text()
		if text == "" {
			dialog.ShowInformation("Предупреждение", "Пустое сообщение", ctx.Win)
			return
		}
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		authCtx := middleware.WithAuth(ctx, ctxWithTimeout)

		_, errSendMessage := clientGRPC.SendMessage(authCtx, &pb.ChatMessage{
			User: &pb.User{Id: ctx.Prefs.UserID(), Username: ctx.Prefs.Username()},
			Text: text,
		})
		if errSendMessage != nil {
			dialog.ShowError(errSendMessage, ctx.Win)
			return
		}
		msgInput.SetText("")
	}

	updateEnableSendBtn(sendBtn, msgInput.Text())

	msgInput.OnChanged = func(text string) {
		updateEnableSendBtn(sendBtn, text)
	}

	updateSendBtn(ctx, sendBtn)

	ctx.App.Preferences().AddChangeListener(func() {
		updateSendBtn(ctx, sendBtn)
	})

	ctx.AddCloseHandler(func() {
		ctx.Prefs.SetLastInput(msgInput.Text())
	})

	msgInput.OnSubmitted = func(text string) {
		sendBtn.OnTapped()
	}

	return container.NewBorder(nil, nil, nil, sendBtn, msgInput)
}

func updateEnableSendBtn(btn *widget.Button, text string) {
	if text != "" {
		btn.Enable()
	} else {
		btn.Disable()
	}
}

func updateSendBtn(ctx *app.Context, sendBtn *widget.Button) {
	if ctx.Prefs.IsAuth() {
		sendBtn.Show()
	} else {
		sendBtn.Hide()
	}
}
