package neptune

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
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
	AppIn     *App
}

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

func GetKeys(key string) *fyne.StaticResource {
	data := sdata.KeyList[key]
	return data
}

func (Ctx *NewContextPlayer) PlaySound(code string) error {
	var f []byte
	var err error
	if Ctx.AppIn.Config.DefaultSound != "nk-cream" {
		file := fmt.Sprintf("%s/%s",Ctx.AppIn.Config.FSounds[Ctx.AppIn.Config.DefaultSound], code)
		f, err = os.ReadFile(file)
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
			panic("Player.Close failed: " + err.Error())
		}
	}()
	return nil
}
