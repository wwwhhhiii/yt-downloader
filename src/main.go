package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

func main() {
	a := app.New()
	window := a.NewWindow("yt-downloader")
	window.Resize(fyne.NewSize(800, 600))
	settings, err := LaodAppSettings()
	if err != nil {
		window.Show()
		dial := dialog.NewInformation("Settings load error", fmt.Sprintf("%s", err), window)
		dial.Show()
		dial.SetOnClosed(func() { log.Fatalf("settings load error: %s", err) })
	}
	window.SetContent(CreateWindowContent(&window, settings))
	window.ShowAndRun()
}
