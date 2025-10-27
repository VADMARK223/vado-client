package hello

import (
	"context"
	"time"
	"vado-client/api/pb/hello"
	"vado-client/internal/app"
	"vado-client/internal/grpc/client"
	"vado-client/internal/grpc/middleware"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewHelloBox(ctx *app.Context) *fyne.Container {
	input := widget.NewEntry()
	input.SetPlaceHolder("Введите имя")

	label := widget.NewLabel("Пусто...")
	button := widget.NewButton("Поздороваться", func() {
		sendHello(ctx, label, input)
	})

	return container.NewVBox(
		input,
		button,
		label,
	)
}

func sendHello(appCtx *app.Context, label *widget.Label, input *widget.Entry) {
	clientGRPC := hello.NewHelloServiceClient(appCtx.GRPC)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	token := client.GetToken(appCtx.App)
	appCtx.Log.Debugf("Send with token: %s", token)
	authCtx := middleware.WithAuth(appCtx, ctx)
	resp, err := clientGRPC.SayHello(authCtx, &hello.HelloRequest{Name: input.Text})
	if err != nil {
		dialog.ShowInformation("Ошибка токена", err.Error(), appCtx.Win)
		return
	}

	label.SetText(resp.Message)
}
