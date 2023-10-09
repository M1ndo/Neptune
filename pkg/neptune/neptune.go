package neptune

import (
	"os"
	"strconv"

	loggdb "github.com/m1ndo/LogGdb"
)

type App struct {
	Config        *SConfig
	ContextPlayer *NewContextPlayer
	KeyEvent      *KeyEvent
	Logger        *loggdb.Logger
	Running       bool
	CfgVars       *SelectVars
	Sys           *Sys
}

func NewApp() *App {
	a := &App{
		Config:        &SConfig{SoundDir: "", DefaultSound: "nk-cream"},
		KeyEvent:      &KeyEvent{},
		ContextPlayer: &NewContextPlayer{Cache: make(map[string][]byte)},
		Logger:        &loggdb.Logger{},
		CfgVars:       &SelectVars{},
		Running:       false,
		Sys:           &Sys{},
	}
	// Set Logger
	a.Logger = setLogger()
	// Cli args
	a.CfgVars = ParseFlags()
	a.Config.SoundDir = a.CfgVars.Sounddir
	// Find Sounds
	a.Config.AppIn = a
	if err := a.Config.FindSounds(); err != nil {
		a.Logger.Log.Error(err)
	}
	// Set Args
	switch {
	case a.CfgVars.Soundkey != "":
		a.Config.DefaultSound = a.CfgVars.Soundkey
	case a.CfgVars.Download:
		_, errCh := DownloadSounds()
		for err := range errCh {
			a.Logger.Log.Error(err)
		}
	case a.CfgVars.ListSounds:
		Fsounds := a.FoundSounds()
		PrintTableWithAliens(Fsounds)
		os.Exit(0)
	}
	// Read soundkey config
	if err := a.Config.ReadConfig(); err != nil {
		a.Logger.Log.Fatal(err)
	}
	// New Context Player
	a.ContextPlayer.AppIn = a
	if err := a.ContextPlayer.NewContext(); err != nil {
		a.Logger.Log.Fatal(err)
	}
	// Sys
	a.Sys.App = a
	return a
}

func (a *App) AppRun() {
	if a.Running {
		a.Logger.Log.Warn("App is already running")
		return
	}
	a.Running = true
	evChan := a.KeyEvent.StartEven()
	a.Logger.Log.Infof("Using Selected Sound %s, ", a.Config.DefaultSound)
	for ev := range evChan {
		go func(event Event) {
			if event.KindCode() == 5 {
				code := strconv.Itoa(int(event.Keycode()))
				if a.Config.Config.KSound[code] != "" {
					if err := a.ContextPlayer.PlaySound(a.Config.Config.KSound[code]); err != nil {
						a.Logger.Log.Error(err)
					}
				}
			}
		}(&CustomEvent{
			Kind:     ev.Kind,
			Keycodes: ev.Keycode,
		})
	}
}

func (a *App) AppStop() {
	if !a.Running {
		a.Logger.Log.Warn("Can't stop and app that is not running")
		return
	}
	a.KeyEvent.StopEvent()
	a.Running = false
}

// List Found Sounds.
func (a *App) FoundSounds() []string {
	var Fsounds []string
	Fsounds = append(Fsounds, "nk-cream")
	for Found := range a.Config.FSounds {
		if Found != "nk-cream" {
			Fsounds = append(Fsounds, Found)
		}
	}
	return Fsounds
}

// Set Soundkeys
func (a *App) SetSounds(sound string) {
	a.Logger.Log.Infof("Now using keysounds %s", sound)
	a.Config.DefaultSound = sound
	if err := a.Config.ReadConfig(); err != nil {
		a.Logger.Log.Fatal(err)
	}
	a.ContextPlayer.ClearCache()
}

// Set Logger Options
func setLogger() *loggdb.Logger {
	dir, _ := GetUserSoundDir()
	NewLogger := &loggdb.Logger{}
	NewLogger.LogDir = dir + "/log"
	CustomOpt := &loggdb.CustomOpt{
		Prefix:          "Neptune 👾 ",
		ReportTimestamp: true,
	}
	NewLogger.LogOptions = CustomOpt
	NewLogger.NewLogger()
	return NewLogger
}
