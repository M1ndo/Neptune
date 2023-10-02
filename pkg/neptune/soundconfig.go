package neptune

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

type SConfig struct {
	SoundDir     string
	DefaultSound string
	Config       *Config
	FSounds      map[string]string
	AppIn        *App
}

type Config struct {
	Name   string            `json:"name"`
	KSound map[string]string `json:"defines"`
}

func (c *SConfig) FindSounds() {
	c.FSounds = make(map[string]string)
	if c.SoundDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			c.AppIn.Logger.Log.Error(err)
		}
		c.SoundDir = filepath.Join(homeDir, "Downloads", "Neptune")
	}
	err := filepath.Walk(c.SoundDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == c.SoundDir {
			return nil
		}
		if info.IsDir() {
			c.AppIn.Logger.Log.Infof("Found Sound! %s", info.Name())
			c.FSounds[info.Name()] = path
		}
		return nil
	})
	if err != nil {
		c.AppIn.Logger.Log.Debug(err)
	}
}

// Open and read config file
func (c *SConfig) ReadConfig() error {
	if c.Config == nil {
		c.Config = &Config{}
	}
	switch {
	case c.DefaultSound == "nk-cream":
		configFile := GetKeys("NkCreamConfig").StaticContent
		fopen := bytes.NewReader(configFile)
		data := json.NewDecoder(fopen)
		if err := data.Decode(c.Config); err != nil {
			return err
		}
	default:
		file := c.FSounds[c.DefaultSound] + "/config.json"
		fopen, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fopen.Close()
		data := json.NewDecoder(fopen)
		if err := data.Decode(c.Config); err != nil {
			return err
		}
	}
	return nil
}
