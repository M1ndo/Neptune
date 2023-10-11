package neptune

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/m1ndo/Neptune/pkg/sdata"
)

type NewContextPlayer struct {
	Options    *oto.NewContextOptions
	Player     *oto.Player
	Context    *oto.Context
	readyChan  chan struct{}
	Cache      map[string][]byte
	mutex      sync.Mutex
	rwmutex    sync.RWMutex
	keyrelease sync.Map
	volume     float64
	AppIn      *App
}

// Decode the file.
func decoder(fBytes *bytes.Reader, ext string) (io.Reader, error) {
	switch {
	case ext == ".wav":
		decoder, err := wav.DecodeWithoutResampling(fBytes)
		return decoder, err
	case ext == ".ogg":
		decoder, err := vorbis.DecodeWithoutResampling(fBytes)
		return decoder, err
	}
	return nil, nil
}

// Create new ContextPlayer
func (Ctx *NewContextPlayer) NewContext() error {
	Ctx.Options = &oto.NewContextOptions{
		SampleRate:   48000,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}
	var err error
	Ctx.Context, Ctx.readyChan, err = oto.NewContext(Ctx.Options)
	if err != nil {
		return err
	}
	return nil
}

// Get Keys
func GetKeys(key string) *fyne.StaticResource {
	data, ok := sdata.KeyList[key]
	if !ok {
		key := string(rune(rand.Intn(26)+'a')) + ".wav"
		return sdata.KeyList[key]
	}
	return data
}

// Store new pressed keys in memory.
func (Ctx *NewContextPlayer) readCache(code, file string) ([]byte, error) {
	Ctx.rwmutex.RLock()
	soundfile, exists := Ctx.Cache[code]
	Ctx.rwmutex.RUnlock()
	if exists {
		Ctx.AppIn.Logger.Log.Debugf("Soundfile %s exists in memory cache", file)
		return soundfile, nil
	}

	Ctx.mutex.Lock()
	defer Ctx.mutex.Unlock()

	soundFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	Ctx.rwmutex.Lock()
	defer Ctx.rwmutex.Unlock()
	Ctx.Cache[code] = soundFile
	Ctx.AppIn.Logger.Log.Debugf("code %s memcache++", code)
	return soundFile, nil
}

// Clear Cache
func (Ctx *NewContextPlayer) ClearCache() {
	for key := range Ctx.Cache {
		delete(Ctx.Cache, key)
	}
}

// Adds file suffix to it.
func addSuffixToFileName(file string, suffix string) string {
	ext := filepath.Ext(file)
	base := strings.TrimSuffix(file, ext)
	return base + suffix + ext
}

// Check if array contains a string
func contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

// Find fallback files.
func findFallBackFiles(dir string) ([]string, error) {
	extensions := []string{".wav", ".ogg"}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fallbackPaths []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		ext := filepath.Ext(filename)

		// Check if the file has the correct prefix and extension
		if strings.HasPrefix(filename, "fallback") && contains(extensions, ext) {
			fallbackPath := filepath.Join(dir, filename)
			fallbackPaths = append(fallbackPaths, fallbackPath)
		}
	}

	return fallbackPaths, nil
}

// Find file or return a fileback
func (Ctx *NewContextPlayer) FindFileWithFallback(dir, name string, event bool) (string, error) {
	// Search for files matching the provided name
	matches, err := filepath.Glob(filepath.Join(dir, name+".*"))
	if err != nil {
		return "", err
	}
	// If a matching file is found, modify the file name if event is true
	if len(matches) > 0 {
		file := matches[0]
		if event {
			file = addSuffixToFileName(file, "-up")
		}
		return file, nil
	}
	// Fallback file names
	fallbackFiles, err := findFallBackFiles(dir)
	if err != nil {
		return "", fmt.Errorf("no matching files found and no fallback files available")
	}
	rIndex := rand.Intn(len(fallbackFiles))
	randomFile := fallbackFiles[rIndex]
	if event {
		randomFile = addSuffixToFileName(randomFile, "-up")
	}
	return randomFile, nil
}

// Play Sounds
func (Ctx *NewContextPlayer) PlaySound(code string, event bool) error {
	var f []byte
	var err error
	var file string
	if Ctx.AppIn.Config.DefaultSound != "nk-cream" {
		if file, err = Ctx.FindFileWithFallback(Ctx.AppIn.Config.FSounds[Ctx.AppIn.Config.DefaultSound], code, event); err != nil {
			return err
		}
		Ctx.AppIn.Logger.Log.Debugf("Playing file %s", file)
		if event {
			code = code + "-up"
		}
		f, err = Ctx.readCache(code, file)
		if err != nil {
			return err
		}
	} else {
		file = fmt.Sprintf("%s.wav", code)
		Ctx.AppIn.Logger.Log.Debugf("Playing Nk-Cream Key %s", file)
		f = GetKeys(file).StaticContent
	}
	fBytes := bytes.NewReader(f)
	decoded, err := decoder(fBytes, path.Ext(file))
	if err != nil {
		return err
	}
	<-Ctx.readyChan
	Ctx.Player = Ctx.Context.NewPlayer(decoded)
	Ctx.Player.SetVolume(Ctx.volume)
	go func() {
		Ctx.Player.Play()
		for Ctx.Player.IsPlaying() {
			time.Sleep(time.Millisecond)
		}
		err := Ctx.Player.Close()
		if err != nil {
			Ctx.AppIn.Logger.Log.Error(err)
		}
	}()
	return nil
}
