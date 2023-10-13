package main

import (
	"github.com/m1ndo/Neptune/pkg/neptune"
	"github.com/m1ndo/Neptune/pkg/ui"
)

func main() {
	App := neptune.NewApp()
	if !App.CfgVars.Cli {
		uiApp := &ui.UiApp{AppIn: App}
		err := uiApp.NewApp()
		if err != nil {
			App.Logger.Log.Fatal(err)
		}
		go uiApp.SystrayRun()
		uiApp.MainWindow.ShowAndRun()
	} else {
		App.Sys.Init()
	}
}
