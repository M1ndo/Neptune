package neptune

import (
	"math/rand"
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
	handleArguments(a)
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

// App Run
func (a *App) AppRun() {
	if a.Running {
		a.Logger.Log.Warn("App is already running")
		return
	}
	a.Running = true
	evChan := a.KeyEvent.StartEven()
	a.Logger.Log.Infof("Using Selected Sound %s, ", a.Config.DefaultSound)
	for ev := range evChan {
		go a.handleEvent(&CustomEvent{
			Kind:     ev.Kind,
			Keycodes: ev.Rawcode,
		})
	}
}

// Handle events
func (a *App) handleEvent(event *CustomEvent) {
	keycode := GetKeyCode(event.Keycodes)
	code := strconv.Itoa(int(keycode))

	if event.Kind == KeyDown {
		a.playSound(code, true)
	} else if event.Kind == KeyUp {
		a.playSound(code, false)
	}
}

// Play
func (a *App) playSound(code string, event bool) {
	if sound := a.Config.Config.KSound[code]; sound != "" {
		a.Logger.Log.Infof("Playing sound %s", sound)
		if err := a.ContextPlayer.PlaySound(sound, event); err != nil {
			a.Logger.Log.Error(err)
		}
	}
}

// App Stop
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

// Random Keysound
func (a *App) AppRand() {
	Sounds := a.FoundSounds()
	randSound := Sounds[rand.Intn(len(Sounds))]
	a.SetSounds(randSound)
}

// Handle arguments
func handleArguments(a *App) {
	if a.CfgVars.Soundkey != "" {
		a.Config.DefaultSound = a.CfgVars.Soundkey
	}
	if a.CfgVars.Download {
		_, errCh := DownloadSounds()
		for err := range errCh {
			a.Logger.Log.Error(err)
		}
	}
	if a.CfgVars.ListSounds {
		Fsounds := a.FoundSounds()
		PrintTableWithAliens(Fsounds)
		os.Exit(0)
	}
	if a.CfgVars.Verbose {
		a.Logger.Log.SetLevel(loggdb.Debug)
	}
}
