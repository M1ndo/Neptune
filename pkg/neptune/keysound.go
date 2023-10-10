package neptune

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/m1ndo/Neptune/pkg/sdata"
)

type NewContextPlayer struct {
	Options   *oto.NewContextOptions
	Player    *oto.Player
	Context   *oto.Context
	readyChan chan struct{}
	Cache     map[string][]byte
	mutex     sync.Mutex
	rwmutex   sync.RWMutex
	AppIn     *App
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
	data := sdata.KeyList[key]
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

// Play Sounds
func (Ctx *NewContextPlayer) PlaySound(code string, event bool) error {
	var f []byte
	var err error
	if Ctx.AppIn.Config.DefaultSound != "nk-cream" {
		file := fmt.Sprintf("%s/%s", Ctx.AppIn.Config.FSounds[Ctx.AppIn.Config.DefaultSound], code)
		f, err = Ctx.readCache(code, file)
		if err != nil {
			return err
		}
	} else {
		f = GetKeys(code).StaticContent
	}
	fBytes := bytes.NewReader(f)
	decoded, err := decoder(fBytes, path.Ext(code))
	if err != nil {
		return err
	}
	<-Ctx.readyChan
	Ctx.Player = Ctx.Context.NewPlayer(decoded)
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
