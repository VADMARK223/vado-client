package chatTab

import (
	"context"
	"time"
	pb "vado-client/api/pb/chat"
	"vado-client/internal/app"
	"vado-client/internal/grpc/middleware"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func newInputBox(ctx *app.Context, clientGRPC pb.ChatServiceClient) (*widget.Entry, *widget.Button) {
	msgInput := widget.NewEntry()
	msgInput.SetText(ctx.Prefs.LastInput())
	msgInput.SetPlaceHolder("Сообщение...")

	sendBtn := createSendBtn(ctx, msgInput, clientGRPC)
	updateEnableSendBtn(sendBtn, msgInput.Text)
	msgInput.OnChanged = func(text string) {
		updateEnableSendBtn(sendBtn, text)
	}
	updateSendBtn(ctx, sendBtn)

	ctx.App.Preferences().AddChangeListener(func() {
		updateSendBtn(ctx, sendBtn)
	})

	ctx.AddCloseHandler(func() {
		ctx.Prefs.SetLastInput(msgInput.Text)
	})

	msgInput.OnSubmitted = func(text string) {
		sendBtn.OnTapped()
	}

	return msgInput, sendBtn
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

func createSendBtn(appCtx *app.Context, input *widget.Entry, grpc pb.ChatServiceClient) *widget.Button {
	result := widget.NewButton("Отправить", func() {
		text := input.Text
		if text == "" {
			dialog.ShowInformation("Предупреждение", "Пустое сообщение", appCtx.Win)
			return
		}
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		authCtx := middleware.WithAuth(appCtx, ctxWithTimeout)

		_, errSendMessage := grpc.SendMessage(authCtx, &pb.ChatMessage{
			User: &pb.User{Id: appCtx.Prefs.UserID(), Username: appCtx.Prefs.Username()},
			Text: text,
		})
		if errSendMessage != nil {
			dialog.ShowError(errSendMessage, appCtx.Win)
			return
		}
		input.SetText("")
	})

	return result
}
