//go:build !systray
// +build !systray

package ui

func (ui *UiApp) SystrayRun() {
	ui.Logger.Log.Warn("App is not built with a systray.")
}
