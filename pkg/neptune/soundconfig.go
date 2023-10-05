package neptune

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
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

func GetUserSoundDir() (string, error) {
	var SoundDir string
	switch runtime.GOOS {
	case "windows":
		appDataDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		SoundDir = filepath.Join(appDataDir, "Neptune")
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		SoundDir = filepath.Join(homeDir, "Library", "Application Support", "Neptune")
	default: // Linux and other Unix-like systems
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		SoundDir = filepath.Join(homeDir, ".local", "share", "Neptune")
	}
	err := os.MkdirAll(SoundDir, 0755)
	if err != nil {
		return "", err
	}
	return SoundDir, nil
}

// Find sounds in a directory.
func (c *SConfig) FindSounds() error {
	c.FSounds = make(map[string]string)
	if c.SoundDir == "" {
		dir, err := GetUserSoundDir()
		if err != nil {
			return err
		}
		c.SoundDir = dir
	}
	err := filepath.Walk(c.SoundDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == c.SoundDir {
			return nil
		}
		if info.IsDir() {
			configFile := filepath.Join(path, "config.json")
			_, err := os.Stat(configFile)
			if err == nil {
				c.AppIn.Logger.Log.Infof("Found Sound! %s", info.Name())
				c.FSounds[info.Name()] = path
			}
		}
		return nil
	})
	if err != nil {
		c.AppIn.Logger.Log.Debug(err)
	}
	return err
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
		file := filepath.Join(c.FSounds[c.DefaultSound], "/config.json")
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
