package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	loggdb "github.com/m1ndo/LogGdb"
)

type UiApp struct {
	AppIn      NeptuneInterface
	MainWindow fyne.Window
	App        fyne.App
	SoundL     *widget.Select
	NotifMsg   *fyne.Notification
	Logger *loggdb.Logger
}

type NeptuneInterface interface {
	AppRun()
	AppStop()
	AppRand()
	SetSounds(string)
	FoundSounds() []string
	DownloadSounds() (string, chan error)
	Checklock() bool
	SetLogger() *loggdb.Logger
}
