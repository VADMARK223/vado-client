package main

import (
	"context"
	"fmt"
	"time"
	"vado-client/internal/pb/hello"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	a := app.NewWithID("vado-client")
	w := a.NewWindow("Client App")

	input := widget.NewEntry()
	input.SetPlaceHolder("Введите имя")

	label := widget.NewLabel("Ожидается ответ...")
	button := widget.NewButton("Отправить", func() {
		conn, err := grpc.NewClient("localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			label.SetText(fmt.Sprintf("Ошибка соединения: %v", err))
			return
		}
		defer func(conn *grpc.ClientConn) {
			_ = conn.Close()
		}(conn)

		client := hello.NewHelloServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := client.SeyHello(ctx, &hello.HelloRequest{Name: input.Text})
		if err != nil {
			label.SetText(fmt.Sprintf("Ошибка запроса: %v", err))
			return
		}

		label.SetText(resp.Message)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Введите имя для gRPC-запроса:"),
		input,
		button,
		label,
	))

	w.Resize(fyne.NewSize(400, 200))
	w.ShowAndRun()
}
