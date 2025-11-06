package hello

import (
	"context"
	"time"
	"vado-client/api/pb/hello"
	"vado-client/internal/app"
	"vado-client/internal/component/common/userInfo"
	"vado-client/internal/grpc/middleware"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewHelloBox(ctx *app.Context) *fyne.Container {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter your name")

	label := widget.NewLabel("")
	button := widget.NewButton("Say hello", func() {
		sendHello(ctx, label, input)
	})

	return container.NewVBox(
		input,
		button,
		label,
	)
}

func sendHello(ctx *app.Context, label *widget.Label, input *widget.Entry) {
	if ctx.Prefs.IsAuth() == false {
		userInfo.ShowLoginDialog(ctx, nil)
		return
	}
	clientGRPC := hello.NewHelloServiceClient(ctx.GRPC)

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	authCtx := middleware.WithAuth(ctx, ctxWithTimeout)
	resp, err := clientGRPC.SayHello(authCtx, &hello.HelloRequest{Name: input.Text})
	if err != nil {
		dialog.ShowInformation("Ошибка токена", err.Error(), ctx.Win)
		return
	}

	label.SetText(resp.Message)
}
