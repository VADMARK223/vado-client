package hello

import (
	"context"
	"fmt"
	"time"
	"vado-client/internal/pb/hello"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewHelloBox() *fyne.Container {
	input := widget.NewEntry()
	input.SetPlaceHolder("Введите имя")

	label := widget.NewLabel("Пусто...")
	button := widget.NewButton("Отправить", func() {
		sendHello(label, input)
	})

	return container.NewVBox(
		input,
		button,
		label,
	)
}

func sendHello(label *widget.Label, input *widget.Entry) {
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
}
