//go:build systray
// +build systray

package ui

import (
	"github.com/getlantern/systray"
	"github.com/m1ndo/Neptune/pkg/sdata"
)

// Register systray
func (ui *UiApp) SystrayRun() {
	systray.Run(ui.OnReady, nil)
}

// onReady() For systray
func (ui *UiApp) OnReady() {
	systray.SetTemplateIcon(sdata.IcoRes.Content(), sdata.IcoRes.Content())
	systray.SetIcon(sdata.IcoRes.Content())
	systray.SetTitle("Neptune")
	systray.SetTooltip("Neptune")
	systray.AddSeparator()
	mShow := systray.AddMenuItem("Show", "Show the main app")
	mStart := systray.AddMenuItem("Start", "Start the soundkeys")
	mPause := systray.AddMenuItem("Stop", "Stop the soundkeys")
	mRand := systray.AddMenuItem("Rand", "Use a random soundkey")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				ui.MainWindow.Show()
			case <-mStart.ClickedCh:
				ui.AppIn.AppRun()
			case <-mPause.ClickedCh:
				ui.AppIn.AppStop()
			case <-mRand.ClickedCh:
				ui.AppIn.AppRand()
			case <-mQuitOrig.ClickedCh:
				ui.MainWindow.Close()
				systray.Quit()
			}
		}
	}()
}
