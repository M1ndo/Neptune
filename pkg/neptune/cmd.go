package neptune

import (
	"flag"
)

type SelectVars struct {
	Cli      bool
	Soundkey string
	Sounddir string
	Download bool
}

func ParseFlags() *SelectVars {
	cfgVars := &SelectVars{}
	soundDir, _ := GetUserSoundDir()
	flag.BoolVar(&cfgVars.Cli, "cli", false, "Run in CLI instead of GUI")
	flag.StringVar(&cfgVars.Soundkey,"soundkey", "", "Soundkey to use default (nk-cream)")
	flag.StringVar(&cfgVars.Sounddir, "sounddir", soundDir, "Sounds directory")
	flag.BoolVar(&cfgVars.Download, "download", false, "Download all other soundkeys")
	flag.Parse()
	return cfgVars
}
