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
	clearBtn  *widget.Button
	container *fyne.Container

	userOnChanged func(string)
	maxLen        int
}

func NewVadoEntry() *VadoEntry {
	e := &VadoEntry{}
	e.ExtendBaseWidget(e)

	e.entry = widget.NewEntry()

	e.clearBtn = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		e.entry.SetText("")
	})
	e.clearBtn.Importance = widget.LowImportance
	e.clearBtn.Disable()

	e.counter = widget.NewLabel("")
	e.counter.Hide()

	e.entry.OnChanged = func(s string) {
		if e.maxLen > 0 {
			runes := []rune(s)
			if len(runes) > e.maxLen {
				e.entry.SetText(string(runes[:e.maxLen]))
			}

			e.updateCounter()
		}

		if s == "" {
			e.clearBtn.Disable()
		} else {
			e.clearBtn.Enable()
		}

		if e.userOnChanged != nil {
			e.userOnChanged(s)
		}
	}

	rightBox := container.NewHBox(e.counter, e.clearBtn)
	alignedRight := container.NewBorder(nil, nil, nil, rightBox, layout.NewSpacer())

	e.container = container.NewStack(e.entry, alignedRight)
	return e
}

func (b *VadoEntry) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(b.container)
}

func (b *VadoEntry) Text() string {
	return b.entry.Text
}

func (b *VadoEntry) SetText(input string) {
	b.entry.SetText(input)
}

func (b *VadoEntry) SetPlaceHolder(s string) {
	b.entry.SetPlaceHolder(s)
}

func (b *VadoEntry) OnChanged(fn func(string)) {
	b.userOnChanged = fn
}

func (b *VadoEntry) OnSubmitted(fn func(string)) {
	b.entry.OnSubmitted = fn
}

func (b *VadoEntry) SetMaxLen(maxLen int) {
	b.maxLen = maxLen
	b.counter.Show()
	b.updateCounter()
}

func (b *VadoEntry) updateCounter() {
	b.counter.SetText(fmt.Sprintf("%d/%d", len([]rune(b.entry.Text)), b.maxLen))
}
