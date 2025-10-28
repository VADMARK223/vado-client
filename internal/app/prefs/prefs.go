package prefs

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
)

type Key string

const (
	LastInput    = "last_input"
	UserID       = "user_id"
	Username     = "username"
	AccessToken  = "access_token" // /home/vadmark/.var/app/com.jetbrains.GoLand/config/fyne/vado-client
	RefreshToken = "refresh_token"
	ExpiresAt    = "expires_at"
)

var allKeys = []Key{
	LastInput,
	UserID,
	Username,
	AccessToken,
	RefreshToken,
	ExpiresAt,
}

type Prefs struct {
	p             fyne.Preferences
	debounceTimer *time.Timer
}

func New(p fyne.Preferences) *Prefs {
	return &Prefs{p: p}
}

func (pr *Prefs) ChangeListeners(callback func()) {
	pr.p.AddChangeListener(func() {
		if pr.debounceTimer != nil {
			pr.debounceTimer.Stop()
		}

		pr.debounceTimer = time.AfterFunc(100*time.Millisecond, func() {
			callback()
		})
	})
}

func (pr *Prefs) UserID() uint64 {
	s := pr.p.String(UserID)
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return u
}

func (pr *Prefs) Reset() {
	for _, k := range allKeys {
		pr.p.RemoveValue(string(k))
	}
}

func (pr *Prefs) SetUserID(value uint64) {
	pr.p.SetString(UserID, strconv.FormatUint(value, 10))
}

func (pr *Prefs) Username() string {
	return pr.p.String(Username)
}

func (pr *Prefs) SetUsername(value string) {
	pr.p.SetString(Username, value)
}

func (pr *Prefs) LastInput() string {
	return pr.p.String(LastInput)
}

func (pr *Prefs) SetLastInput(value string) {
	pr.p.SetString(LastInput, value)
}

func (pr *Prefs) AccessToken() string {
	return pr.p.String(AccessToken)
}

func (pr *Prefs) SetAccessToken(value string) {
	pr.p.SetString(AccessToken, value)
}

func (pr *Prefs) IsAuth() bool {
	return pr.AccessToken() != ""
}

func (pr *Prefs) RefreshToken() string {
	return pr.p.String(RefreshToken)
}

func (pr *Prefs) SetRefreshToken(value string) {
	pr.p.SetString(RefreshToken, value)
}

func (pr *Prefs) ExpiresAt() int64 {
	s := pr.p.String(ExpiresAt)
	u, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return u
}

func (pr *Prefs) SetExpiresAt(value int64) {
	pr.p.SetString(ExpiresAt, strconv.Itoa(int(value)))
}
