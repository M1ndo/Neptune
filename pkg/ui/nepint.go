package ui

type NeptuneInterface interface {
	AppRun()
	AppStop()
	AppRand()
	SetSounds(string)
	FoundSounds() []string
}
