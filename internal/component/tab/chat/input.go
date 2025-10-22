package chat

import (
	"context"
	pb "vado-client/api/pb/chat"
	"vado-client/internal/appcontext"
	"vado-client/internal/grpc/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func newInputBox(appCtx *appcontext.AppContext, ctx context.Context, clientGRPC pb.ChatServiceClient) *fyne.Container {
	msgInput := widget.NewEntry()
	msgInput.SetText(client.GetLastInput(appCtx.App))
	msgInput.SetPlaceHolder("Сообщение...")

	sendBtn := createSendBtn(appCtx, ctx, msgInput, clientGRPC)
	updateEnableSendBtn(sendBtn, msgInput.Text)
	msgInput.OnChanged = func(text string) {
		updateEnableSendBtn(sendBtn, text)
	}
	updateSendBtn(appCtx.App, sendBtn)

	appCtx.App.Preferences().AddChangeListener(func() {
		updateSendBtn(appCtx.App, sendBtn)
	})

	appCtx.AddCloseHandler(func() {
		client.SetLastInput(appCtx.App, msgInput.Text)
	})

	return container.NewVBox(msgInput, sendBtn)
}

func updateEnableSendBtn(btn *widget.Button, text string) {
	if text != "" {
		btn.Enable()
	} else {
		btn.Disable()
	}
}

func updateSendBtn(app fyne.App, sendBtn *widget.Button) {
	if client.IsAuth(app) {
		sendBtn.Show()
	} else {
		sendBtn.Hide()
	}
}

func createSendBtn(appCtx *appcontext.AppContext, ctx context.Context, input *widget.Entry, grpc pb.ChatServiceClient) *widget.Button {
	result := widget.NewButton("Отправить", func() {
		text := input.Text
		if text == "" {
			dialog.ShowInformation("Предупреждение", "Пустое сообщение", appCtx.Win)
			return
		}
		token := client.GetToken(appCtx.App)
		appCtx.Log.Debugf("Send with token: %s", token)
		authCtx := withAuth(ctx, token)

		_, errSendMessage := grpc.SendMessage(authCtx, &pb.ChatMessage{
			Id:   client.GetUserID(appCtx.App),
			User: client.GetUsername(appCtx.App),
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
