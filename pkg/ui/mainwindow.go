package ui

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	// "fyne.io/fyne/v2/storage"
	// xwidget "fyne.io/x/fyne/widget"
)

type UiApp struct {
	MainWindow fyne.Window
}

func (ui *UiApp) NewApp(SoundL fyne.CanvasObject, start func(), stop func()) error {
	app := app.NewWithID("cf.ybenel.Neptune")
	app.Settings().SetTheme(&myTheme{})
	app.SetIcon(IconRes)
	w := app.NewWindow("Neptune")
	w.Resize(fyne.NewSize(460, 400))
	w.SetFixedSize(true)
	w.CenterOnScreen()
	// Create a box container
	// box := container.NewVBox()
	box := container.New(layout.NewVBoxLayout())

	// Create a label with centered text and a purple background
	purple := color.RGBA{R: 128, G: 0, B: 128, A: 255}
	title := canvas.NewText("Neptune", purple)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 24.0

	// More Text
	textString := "Written By ybenel In Golang  💙 Using Emacs 👾"
	authText := widget.NewRichTextFromMarkdown(textString)
	nAuthText := container.NewCenter(authText)

	// Logo
	logo := canvas.NewImageFromResource(NeptuneRes)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(256, 256))
	// Spinning Logo
	// Slogo, err := xwidget.NewAnimatedGif(storage.NewFileURI("neptunes.gif"))
	// Slogo.SetMinSize(fyne.NewSize(256, 256))
	// if err != nil {
	// 	return err
	// }

	// Start And Stop
	StartStop := container.NewGridWithColumns(2,
		&widget.Button{
			Text:       "Start",
			Alignment:  widget.ButtonAlignCenter,
			Importance: widget.HighImportance,
			OnTapped: func() {
				go start()
			},
		},
		&widget.Button{
			Text:       "Stop",
			Alignment:  widget.ButtonAlignCenter,
			Importance: widget.WarningImportance,
			OnTapped:   stop,
		},
	)

	// Links
	NLinks := container.NewHBox(
		widget.NewHyperlink("Neptune", parseURL("https://github.com/m1ndo/Neptune")),
		widget.NewLabel("-"),
		widget.NewHyperlink("How To", parseURL("https://github.com/m1ndo/Neptune")),
		widget.NewLabel("-"),
		widget.NewHyperlink("Donate", parseURL("https://ybenel.cf/donate")),
	)

	// Create a spacer to push the buttons down
	buttonsSpacer := container.NewVBox(widget.NewLabel(""))
	buttonsSpacer.Resize(fyne.NewSize(10, 800))

	// Add Widgets
	box.Add(title)
	// box.Add(Slogo)
	box.Add(logo)
	box.Add(SoundL)
	box.Add(buttonsSpacer)
	box.Add(StartStop)
	// box.Add(buttonsSpacer)
	// box.Add(buttonsSpacer)
	box.Add(nAuthText)
	// box.Add(buttonsSpacer)
	// box.Add(buttonsSpacer)
	box.Add(container.NewCenter(NLinks))
	// Slogo.Start()
	w.SetContent(box)
	ui.MainWindow = w
	return nil
}

func (Ui *UiApp) SoundsList(SoundList []string, f func(string)) fyne.CanvasObject {
	AvailableSounds := widget.NewSelect(SoundList, f)
	AvailableSounds.Selected = "nk-cream"
	return AvailableSounds
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}
