package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kkdai/youtube/v2"
)

// TIP: to show pgbar:
// go func() {
// 	time.Sleep(time.Second * 3)
// 	dwnldProgressBar.Show()
// 	fyne.Do(vbox.Refresh)
// 	fmt.Println("pgbar shown")
// }()

var testVideoId string = "0o8Ex8mXigU"

type DownloadOptions struct {
	saveDirectory string
}

type DownloadResult struct {
	err error
}

func DownloadVideo(videoId string, opts *DownloadOptions, pbar *widget.ProgressBar, vbox *fyne.Container, errChan chan<- DownloadResult) {
	fyne.Do(pbar.Show)
	fyne.Do(vbox.Refresh)

	client := youtube.Client{} // TODO: reuse client?
	video, err := client.GetVideo(videoId)
	if err != nil {
		errChan <- DownloadResult{err: err}
	}
	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		errChan <- DownloadResult{err: err}
	}
	defer stream.Close()

	os.MkdirAll(opts.saveDirectory, os.ModePerm)

	file, err := os.Create(filepath.Join(opts.saveDirectory, fmt.Sprintf("%s.mp4", video.Title)))
	if err != nil {
		errChan <- DownloadResult{err: err}
	}
	defer file.Close()

	fyne.Do(func() { pbar.SetValue(pbar.Max / 2) })
	_, err = io.Copy(file, stream)
	if err != nil {
		errChan <- DownloadResult{err: err}
	}
	fyne.Do(func() { pbar.SetValue(pbar.Max) })
	fyne.Do(pbar.Hide)
	fyne.Do(vbox.Refresh)
	errChan <- DownloadResult{err: nil}
}

func CreateWindowContent(forWindow *fyne.Window, settings *AppSettings) *fyne.Container {
	vbox := container.NewVBox()
	videoEntry := widget.NewEntry()
	dwnldProgressBar := widget.NewProgressBar()
	dwnldProgressBar.Hide()
	dwnldBtn := widget.NewButton(
		"Download",
		func() {
			resultChan := make(chan DownloadResult)
			go DownloadVideo(
				videoEntry.Text,
				&DownloadOptions{saveDirectory: settings.SaveDirectory},
				dwnldProgressBar,
				vbox,
				resultChan,
			)
			if result := <-resultChan; result.err != nil {
				dial := dialog.NewInformation("Request error", result.err.Error(), *forWindow)
				dial.Show()
			}
		},
	)
	videoEntry.SetPlaceHolder("Video/playlist ID/URL")
	vbox.Add(videoEntry)
	vbox.Add(dwnldProgressBar)
	vbox.Add(dwnldBtn)
	return vbox
}
