package app

import (
	"vado-client/internal/app/keyman"

	"fyne.io/fyne/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Context struct {
	Log      *zap.SugaredLogger
	App      fyne.App
	Win      fyne.Window
	GRPC     *grpc.ClientConn
	OnClose  []func()
	KeyMan   *keyman.KeyManager
	unsubEsc func()
}

func NewAppContext(log *zap.SugaredLogger, grpc *grpc.ClientConn, a fyne.App, w fyne.Window) *Context {
	result := &Context{
		Log:    log,
		App:    a,
		Win:    w,
		GRPC:   grpc,
		KeyMan: keyman.New(w.Canvas()),
	}

	result.AddCloseHandler(func() {
		_ = grpc.Close()
	})
	return result
}

func (a *Context) AddCloseHandler(fn func()) {
	a.OnClose = append(a.OnClose, fn)
	a.Log.Debugw("Added close handler", "count", len(a.OnClose))

	a.Win.SetCloseIntercept(func() {
		a.Log.Debugw("Close intercepted.")
		a.callCloseHandlers()
		a.Dispose()
		a.Win.Close()
	})

	if a.unsubEsc == nil {
		a.unsubEsc = a.KeyMan.Subscribe(func(ev *fyne.KeyEvent) {
			if ev.Name == fyne.KeyEscape {
				a.Log.Debugw("Press escape. Close.")
				a.callCloseHandlers()
				a.Win.Close()
			}
		})
	}
}

func (a *Context) callCloseHandlers() {
	for _, f := range a.OnClose {
		f()
	}
}

func (a *Context) Dispose() {
	if a.unsubEsc != nil {
		a.unsubEsc()
		a.unsubEsc = nil
	}
}
