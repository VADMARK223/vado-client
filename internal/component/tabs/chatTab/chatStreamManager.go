package chatTab

import (
	"context"
	"errors"
	"io"
	"time"
	pb "vado-client/api/pb/chat"
	"vado-client/internal/app"
	"vado-client/internal/grpc/middleware"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type ChatStreamManager struct {
	appCtx   *app.Context
	client   pb.ChatServiceClient
	messages binding.UntypedList
	active   bool
	cancel   context.CancelFunc
	onSystem func(count uint32)
}

func NewChatStreamManager(appCtx *app.Context, client pb.ChatServiceClient, messages binding.UntypedList, onSystem func(count uint32)) *ChatStreamManager {
	m := &ChatStreamManager{
		appCtx:   appCtx,
		client:   client,
		messages: messages,
		onSystem: onSystem,
	}
	appCtx.AddCloseHandler(m.Stop)
	return m
}

func (m *ChatStreamManager) Start() {
	if m.active || !m.appCtx.Prefs.IsAuth() {
		return
	}
	m.active = true

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	go func() {
		m.appCtx.Log.Infow("Start chat stream manager.")
		defer func() {
			m.Stop()
			m.appCtx.Log.Infow("Stop chat stream manager.")
		}()

		for {
			req := &pb.ChatStreamRequest{
				User: &pb.User{
					Id:       m.appCtx.Prefs.UserID(),
					Username: m.appCtx.Prefs.Username(),
				},
			}
			stream, errChatStream := m.client.ChatStream(middleware.WithAuth(m.appCtx, ctx), req)
			if errChatStream != nil {
				m.appCtx.Log.Errorw("Error connecting to stream", "error", errChatStream)
				select {
				case <-time.After(2 * time.Second):
					continue
				case <-ctx.Done():
					return
				}
			}

			for {
				msg, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) || ctx.Err() != nil {
						m.appCtx.Log.Debug("stream closed by client")
						return
					}
					m.appCtx.Log.Errorw("Error receiving message", "error", err)
					break
				}

				if err != nil {
					if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) || ctx.Err() != nil {
						m.appCtx.Log.Info("stream closed by client")
						return
					}
					m.appCtx.Log.Errorw("Error receiving message", "error", err)

					break
				}

				fyne.Do(func() {
					// Пропускаем собственное системное сообщение
					if msg.Type != pb.MessageType_MESSAGE_SYSTEM || msg.User.Id != m.appCtx.Prefs.UserID() {
						if err := m.messages.Append(msg); err != nil {
							m.appCtx.Log.Errorw("Error append message", "error", err)
						}
					}

					if msg.Type == pb.MessageType_MESSAGE_SYSTEM {
						m.onSystem(msg.UsersCount)
					}
				})
			}
		}
	}()
}

func (m *ChatStreamManager) Stop() {
	if m.cancel != nil {
		m.appCtx.Log.Infow("ChatStreamManager cancel.")
		m.cancel()
		m.cancel = nil
	}
	m.active = false
}
