package keyman

import (
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
)

type KeyManager struct {
	mu        sync.Mutex
	canvas    fyne.Canvas
	listeners map[uint64]Handler
	nextID    uint64
}

type Handler func(ev *fyne.KeyEvent)

func New(canvas fyne.Canvas) *KeyManager {
	m := &KeyManager{
		canvas:    canvas,
		listeners: make(map[uint64]Handler),
	}

	canvas.SetOnTypedKey(func(ev *fyne.KeyEvent) {
		m.mu.Lock()
		defer m.mu.Unlock()
		for _, h := range m.listeners {
			h(ev)
		}
	})

	return m
}

// Subscribe добавляет обработчик и возвращает функцию для отписки
func (m *KeyManager) Subscribe(h Handler) func() {
	id := atomic.AddUint64(&m.nextID, 1)

	m.mu.Lock()
	m.listeners[id] = h
	m.mu.Unlock()

	// возвращаем функцию для отписки
	return func() {
		m.mu.Lock()
		delete(m.listeners, id)
		m.mu.Unlock()
	}
}
