package common

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type VadoEntry struct {
	widget.BaseWidget
	entry     *widget.Entry
	counter   *widget.Label
	clearBtn  *VadoIconButton
	container *fyne.Container

	OnSubmitted  func(string)
	OutOnChanged func(string)

	maxLen int
}

func NewVadoEntry() *VadoEntry {
	e := &VadoEntry{}
	e.ExtendBaseWidget(e)

	e.entry = widget.NewEntry()
	e.entry.OnSubmitted = func(s string) {
		if e.OnSubmitted != nil {
			e.OnSubmitted(s)
		}
	}
	e.clearBtn = NewVadoIconButton(theme.CancelIcon(), func() {
		e.entry.SetText("")
	})
	e.clearBtn.Disable()

	e.counter = widget.NewLabel("")
	e.counter.Hide()

	e.entry.OnChanged = e.internalOnChanged

	rightBox := container.NewHBox(e.counter, e.clearBtn)
	alignedRight := container.NewBorder(nil, nil, nil, rightBox, layout.NewSpacer())

	e.container = container.NewStack(e.entry, alignedRight)
	return e
}

func (e *VadoEntry) internalOnChanged(s string) {
	if e.maxLen > 0 {
		runes := []rune(s)
		if len(runes) > e.maxLen {
			e.entry.OnChanged = nil
			e.entry.SetText(string(runes[:e.maxLen]))
			e.entry.OnChanged = e.internalOnChanged
			return
		}
		e.updateCounter()
	}

	if s == "" {
		e.clearBtn.Disable()
	} else {
		e.clearBtn.Enable()
	}
	if e.OutOnChanged != nil {
		e.OutOnChanged(s)
	}
}

func (e *VadoEntry) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(e.container)
}

func (e *VadoEntry) Text() string {
	return e.entry.Text
}

func (e *VadoEntry) SetText(input string) {
	e.entry.SetText(input)
}

func (e *VadoEntry) SetPlaceHolder(s string) {
	e.entry.SetPlaceHolder(s)
}

func (e *VadoEntry) SetMaxLen(maxLen int) {
	e.maxLen = maxLen
	e.counter.Show()
	e.updateCounter()
}

func (e *VadoEntry) updateCounter() {
	e.counter.SetText(fmt.Sprintf("%d/%d", len([]rune(e.entry.Text)), e.maxLen))
}
