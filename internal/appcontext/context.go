package appcontext

import (
	"fyne.io/fyne/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AppContext struct {
	Log  *zap.SugaredLogger
	GRPC *grpc.ClientConn
	Win  fyne.Window
}

func NewAppContext(log *zap.SugaredLogger, client *grpc.ClientConn, w fyne.Window) *AppContext {
	return &AppContext{
		Log:  log,
		GRPC: client,
		Win:  w,
	}
}
