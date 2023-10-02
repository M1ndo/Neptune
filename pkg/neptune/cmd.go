package neptune

import (
	"flag"
)

type SelectVars struct {
	Cli      bool
	Soundkey string
}

func ParseFlags() *SelectVars {
	cfgVars := &SelectVars{}
	flag.BoolVar(&cfgVars.Cli, "cli", false, "Run in CLI instead of GUI")
	soundKey := flag.String("soundkey", "", "Soundkey to use default (nk-cream)")
	flag.Parse()
	if *soundKey != "" {
		cfgVars.Soundkey = *soundKey
	}
	return cfgVars
}
