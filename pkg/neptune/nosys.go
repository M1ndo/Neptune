//go:build !systray
// +build !systray

package neptune

type Sys struct {
	*App
}

func (App *Sys) Init() {
	App.Logger.Log.Warn("App is not built with a systray")
	App.AppRun()
}
