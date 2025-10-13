package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	pb "vado-client/internal/pb/chat"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	a := app.NewWithID("vado-client")
	w := a.NewWindow("Vado client")
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			closeWindow(w)
		}
	})

	// подключаемся к gRPC серверу
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("не удалось подключиться: %v", err)
	}

	client := pb.NewChatServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())

	// При закрытии окна — корректно завершаем соединение
	w.SetCloseIntercept(func() {
		cancel()
		conn.Close()
		a.Quit()
	})

	messages := widget.NewMultiLineEntry()
	messages.Disable()

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Ваше имя")

	input := widget.NewEntry()
	input.SetPlaceHolder("Сообщение...")

	sendButton := widget.NewButton("Отправить", func() {
		user := nameEntry.Text
		text := input.Text
		if user == "" || text == "" {
			return
		}

		_, err := client.SendMessage(ctx, &pb.ChatMessage{
			User: user,
			Text: text,
		})
		if err != nil {
			log.Printf("Ошибка отправки: %v", err)
			return
		}
		input.SetText("")
	})

	content := container.NewBorder(
		container.NewVBox(nameEntry, messages),
		container.NewVBox(input, sendButton),
		nil, nil,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 600))

	// поток сообщений
	go func() {
		stream, err := client.ChatStream(ctx, &pb.Empty{})
		if err != nil {
			log.Printf("ошибка ChatStream: %v", err)
			return
		}

		for {
			msg, err := stream.Recv()
			if err == io.EOF || ctx.Err() != nil {
				break
			}
			if err != nil {
				log.Printf("ошибка получения: %v", err)
				break
			}
			text := fmt.Sprintf("%s: %s\n", msg.User, msg.Text)
			fyne.Do(func() {
				messages.SetText(messages.Text + text)
			})
		}
	}()

	// Обработка Ctrl+C
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch
		cancel()
		conn.Close()
		a.Quit()
	}()

	w.ShowAndRun()
}

func createChatBox() *fyne.Container {
	conn, errConn := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if errConn != nil {
		return container.NewVBox(widget.NewLabel(fmt.Sprintf("Ошибка соединения: %v", errConn)))
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	client := pb.NewChatServiceClient(conn)
	ctx := context.Background()

	messages := widget.NewMultiLineEntry()
	messages.SetPlaceHolder("Сообщения появятся здесь...")
	messages.Disable()

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите имя")

	input := widget.NewEntry()
	input.SetPlaceHolder("Введите сообщение")

	sendButton := widget.NewButton("Отправить", func() {
		user := nameEntry.Text
		text := input.Text
		if user == "" || text == "" {
			return
		}

		_, err := client.SendMessage(ctx, &pb.ChatMessage{
			User: user,
			Text: text,
		})
		if err != nil {
			log.Printf("Send error: %v", err)
		}
		input.SetText("")
	})

	content := container.NewBorder(
		container.NewVBox(nameEntry, messages),
		container.NewVBox(input, sendButton),
		nil, nil,
	)

	// читаем поток сообщений
	go func() {
		stream, err := client.ChatStream(ctx, &pb.Empty{})
		if err != nil {
			log.Printf("ChatStream error: %v", err)
			return
		}
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Recv error: %v", err)
				break
			}
			//a.Settings().SetTheme(&fyne.CurrentApp().Settings().Theme())
			messages.SetText(fmt.Sprintf("%s[%s]: %s\n%s",
				messages.Text,
				msg.User,
				msg.Text,
				time.Now().Format("15:04:05"),
			))
		}
	}()

	return content
}

func closeWindow(w fyne.Window) {
	w.Close()
}
