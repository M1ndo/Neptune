package main

import (
	"github.com/m1ndo/Neptune/pkg/neptune"
	"github.com/m1ndo/Neptune/pkg/ui"
)

func main() {
	App := neptune.NewApp()
	if !App.CfgVars.Cli {
		uiApp := &ui.UiApp{}
		SoundWid := uiApp.SoundsList(App.FoundSounds(), App.SetSounds)
		err := uiApp.NewApp(SoundWid, App.AppRun, App.AppStop)
		if err != nil {
			App.Logger.Log.Fatal(err)
		}
		uiApp.MainWindow.ShowAndRun()
	} else {
		App.AppRun()
	}
}
