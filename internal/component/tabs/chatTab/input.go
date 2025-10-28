package chatTab

import (
	"context"
	pb "vado-client/api/pb/chat"
	"vado-client/internal/app"
	"vado-client/internal/grpc/middleware"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func newInputBox(appCtx *app.Context, ctx context.Context, clientGRPC pb.ChatServiceClient) (*widget.Entry, *widget.Button) {
	msgInput := widget.NewEntry()
	msgInput.SetText(appCtx.Prefs.LastInput())
	msgInput.SetPlaceHolder("Сообщение...")

	sendBtn := createSendBtn(appCtx, ctx, msgInput, clientGRPC)
	updateEnableSendBtn(sendBtn, msgInput.Text)
	msgInput.OnChanged = func(text string) {
		updateEnableSendBtn(sendBtn, text)
	}
	updateSendBtn(appCtx, sendBtn)

	appCtx.App.Preferences().AddChangeListener(func() {
		updateSendBtn(appCtx, sendBtn)
	})

	appCtx.AddCloseHandler(func() {
		appCtx.Prefs.SetLastInput(msgInput.Text)
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

func createSendBtn(appCtx *app.Context, ctx context.Context, input *widget.Entry, grpc pb.ChatServiceClient) *widget.Button {
	result := widget.NewButton("Отправить", func() {
		text := input.Text
		if text == "" {
			dialog.ShowInformation("Предупреждение", "Пустое сообщение", appCtx.Win)
			return
		}
		authCtx := middleware.WithAuth(appCtx, ctx)

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
