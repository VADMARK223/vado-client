package hello

import (
	"context"
	"time"
	"vado-client/internal/appcontext"
	"vado-client/internal/pb/hello"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewHelloBox(ctx *appcontext.AppContext) *fyne.Container {
	input := widget.NewEntry()
	input.SetPlaceHolder("Введите имя")

	label := widget.NewLabel("Пусто...")
	button := widget.NewButton("Отправить", func() {
		sendHello(ctx, label, input)
	})

	return container.NewVBox(
		input,
		button,
		label,
	)
}

func sendHello(appCtx *appcontext.AppContext, label *widget.Label, input *widget.Entry) {
	client := hello.NewHelloServiceClient(appCtx.GRPC)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.SeyHello(ctx, &hello.HelloRequest{Name: input.Text})
	if err != nil {
		dialog.ShowError(err, appCtx.Win)
		return
	}

	label.SetText(resp.Message)
}
