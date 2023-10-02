package neptune

import (
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
}

func NewApp() *App {
	a := &App{
		Config:        &SConfig{SoundDir: "sounds/", DefaultSound: "nk-cream"},
		KeyEvent:      &KeyEvent{},
		ContextPlayer: &NewContextPlayer{},
		Logger:        &loggdb.Logger{},
		CfgVars:       &SelectVars{},
		Running:       false,
	}
	// Set Logger
	a.Logger.NewLogger()
	// Cli args
	a.CfgVars = ParseFlags()
	if a.CfgVars.Soundkey != "" {
		a.Config.DefaultSound = a.CfgVars.Soundkey
	}
	// Find Sounds
	a.Config.AppIn = a
	a.Config.FindSounds()
	if err := a.Config.ReadConfig(); err != nil {
		a.Logger.Log.Fatal(err)
	}
	// New Context Player
	a.ContextPlayer.AppIn = a
	if err := a.ContextPlayer.NewContext(); err != nil {
		a.Logger.Log.Fatal(err)
	}
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
					// sfile := a.Config.FSounds[a.Config.DefaultSound] + "/" + a.Config.Config.KSound[code]
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
	for Found := range a.Config.FSounds {
		Fsounds = append(Fsounds, Found)
	}
	return Fsounds
}

// Set Soundkeys
func (a *App) SetSounds(sound string) {
	a.Logger.Log.Infof("Executed and sat %s", sound)
	a.Config.DefaultSound = sound
}
