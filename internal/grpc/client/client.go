package client

import (
	"fmt"
	"vado-client/internal/config/port"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateClient() (*grpc.ClientConn, error) {
	target := "localhost:" + port.GRPC
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания gRPC клиента: %w", err)
	}

	return conn, nil
}
