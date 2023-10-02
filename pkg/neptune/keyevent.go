package neptune

import (
	hook "github.com/robotn/gohook"
)

type KeyEvent struct {
	Enabled bool
	EvChan  chan hook.Event
}

type Event interface {
	KindCode() uint8
	Keycode() uint16
}

type CustomEvent struct {
	Kind    uint8
	Keycodes uint16
}


func NewKeyEvent() *KeyEvent {
	return &KeyEvent{}
}

func (Kv *KeyEvent) StartEven() chan hook.Event {
	Kv.Enabled = true
	Kv.EvChan = hook.Start()
	return Kv.EvChan
}

func (Kv *KeyEvent) StopEvent() {
	if Kv.Enabled {
		hook.End()
		Kv.Enabled = false
	}
}

func (k *KeyEvent) EventChannel() chan hook.Event {
	return k.EvChan
}

func (ce CustomEvent) KindCode() uint8 {
	return ce.Kind
}

func (ce CustomEvent) Keycode() uint16 {
	return ce.Keycodes
}
