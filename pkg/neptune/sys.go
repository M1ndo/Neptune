package neptune

import (
	"github.com/getlantern/systray"
	"github.com/m1ndo/Neptune/pkg/ui"
)

type Sys struct {
	*App
}

func (App *Sys) Init() {
	systray.Run(App.onReady, nil)
}

func (App *Sys) onReady() {
	go App.AppRun()
	systray.SetTemplateIcon(ui.IconRes.Content(), ui.IconRes.Content())
	systray.SetIcon(ui.IconRes.Content())
	systray.SetTitle("Neptune")
	systray.SetTooltip("Neptune")
	systray.AddSeparator()
	mStart := systray.AddMenuItem("Start", "Start the soundkeys")
	mPause := systray.AddMenuItem("Stop", "Stop the soundkeys")
	mRand := systray.AddMenuItem("Rand", "Use a random soundkey")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for {
			select {
			case <-mStart.ClickedCh:
				App.AppRun()
			case <-mPause.ClickedCh:
				App.AppStop()
			case <-mRand.ClickedCh:
				App.AppRand()
			case <-mQuitOrig.ClickedCh:
				systray.Quit()
			}
		}
	}()
}
