package neptune

import (
	"flag"
)

type SelectVars struct {
	Cli        bool
	Soundkey   string
	Sounddir   string
	Download   bool
	ListSounds bool
	Verbose    bool
}

func ParseFlags() *SelectVars {
	cfgVars := &SelectVars{}
	soundDir, _ := GetUserSoundDir()
	flag.BoolVar(&cfgVars.Cli, "cli", false, "Run in CLI instead of GUI")
	flag.StringVar(&cfgVars.Soundkey, "soundkey", "", "Soundkey to use default (nk-cream)")
	flag.StringVar(&cfgVars.Sounddir, "sounddir", soundDir, "Sounds directory")
	flag.BoolVar(&cfgVars.Download, "download", false, "Download all other soundkeys")
	flag.BoolVar(&cfgVars.ListSounds, "lst", false, "List all available sounds")
	flag.BoolVar(&cfgVars.Verbose, "verbose", false, "Verbose output (Debugging)")
	flag.Parse()
	return cfgVars
}
