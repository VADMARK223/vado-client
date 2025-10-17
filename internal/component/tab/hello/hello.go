package hello

import (
	"context"
	"time"
	"vado-client/api/pb/hello"
	"vado-client/internal/appcontext"
	"vado-client/internal/grpc/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc/metadata"
)

func NewHelloBox(ctx *appcontext.AppContext, a fyne.App) *fyne.Container {
	input := widget.NewEntry()
	input.SetPlaceHolder("Введите имя")

	label := widget.NewLabel("Пусто...")
	button := widget.NewButton("Поздороваться", func() {
		sendHello(ctx, label, input, client.GetToken(a))
	})

	return container.NewVBox(
		input,
		button,
		label,
	)
}

func withAuth(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

func sendHello(appCtx *appcontext.AppContext, label *widget.Label, input *widget.Entry, token string) {
	client := hello.NewHelloServiceClient(appCtx.GRPC)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	appCtx.Log.Debugf("Send with token: %s", token)
	authCtx := withAuth(ctx, token)
	resp, err := client.SeyHello(authCtx, &hello.HelloRequest{Name: input.Text})
	if err != nil {
		dialog.ShowInformation("Ошибка токена", err.Error(), appCtx.Win)
		return
	}

	label.SetText(resp.Message)
}
