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
	a.Logger = a.SetLogger()
	// Cli args
	a.CfgVars = ParseFlags()
	a.Config.SoundDir = a.CfgVars.Sounddir
	a.Config.AppIn = a
	// Set Args
	handleArguments(a)
	// Find Sounds
	if err := a.Config.FindSounds(); err != nil {
		a.Logger.Log.Error(err)
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
	rawcode := strconv.Itoa(int(keycode))
	code := CodeToChar(keycode)
	if a.Config.Config.IsMulti {
		switch event.Kind {
		case KeyDown:
			_, loaded := a.ContextPlayer.keyrelease.LoadOrStore(keycode, true)
			if loaded {
				return
			}
			a.Logger.Log.Debugf("Clicked down %s keycode %s, Keychar %s", code, rawcode, a.KeyEvent.CodeToChar(keycode))
			a.playSound(code, false)
    case KeyUp:
			a.ContextPlayer.keyrelease.Delete(keycode)
			a.Logger.Log.Debugf("Clicked up %s keycode %s, Keychar %s", code, rawcode, a.KeyEvent.CodeToChar(keycode))
			a.playSound(code, true)
		}
	} else {
		if event.Kind == KeyUp {
			a.Logger.Log.Debugf("Playing keycode %s, Keychar %s", rawcode, a.KeyEvent.CodeToChar(keycode))
			a.playSound(code, false)
		}
	}
}

// Play
func (a *App) playSound(code string, event bool) {
	if code != "" {
		if err := a.ContextPlayer.PlaySound(code, event); err != nil {
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
	if a.Config.DefaultSound == "nk-cream" {
		a.Config.Config.IsMulti = false
	}
}

// Set Logger Options
func (a *App) SetLogger() *loggdb.Logger {
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

// Download Sounds
func (a *App) DownloadSounds() (string, chan error){
	msg, err := DownloadSounds(false)
	return msg, err
}

// Download Sounds
func (a *App) Checklock() bool {
	return checkLock()
}

// Handle arguments
func handleArguments(a *App) {
	if a.CfgVars.Soundkey != "" {
		a.Config.DefaultSound = a.CfgVars.Soundkey
	}
	if a.CfgVars.Download {
		msg, errCh := DownloadSounds(true)
		if errCh != nil {
			for err := range errCh {
				a.Logger.Log.Error(err)
			}
		}
		if msg != "" {
			a.Logger.Log.Info(msg)
		}
	}
	if a.CfgVars.ListSounds {
		if err := a.Config.FindSounds(); err != nil {
			a.Logger.Log.Error(err)
		}
		Fsounds := a.FoundSounds()
		PrintTableWithAliens(Fsounds)
		os.Exit(0)
	}
	if a.CfgVars.Volume != 0 {
		a.ContextPlayer.volume = a.CfgVars.Volume
	}
	if a.CfgVars.Verbose {
		a.Logger.Log.SetLevel(loggdb.Debug)
	}
}
