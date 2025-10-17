package appcontext

import (
	"fyne.io/fyne/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AppContext struct {
	Log     *zap.SugaredLogger
	App     fyne.App
	Win     fyne.Window
	GRPC    *grpc.ClientConn
	OnClose []func()
}

func NewAppContext(log *zap.SugaredLogger, grpc *grpc.ClientConn, a fyne.App, w fyne.Window) *AppContext {
	result := &AppContext{
		Log:  log,
		App:  a,
		Win:  w,
		GRPC: grpc,
	}

	result.AddCloseHandler(func() {
		_ = grpc.Close()
	})
	return result
}

func (a *AppContext) AddCloseHandler(fn func()) {
	a.OnClose = append(a.OnClose, fn)
	a.Win.SetCloseIntercept(func() {
		for _, f := range a.OnClose {
			f()
		}
		a.Win.Close()
	})

	a.Win.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			for _, f := range a.OnClose {
				f()
			}
			a.Win.Close()
		}
	})
}
