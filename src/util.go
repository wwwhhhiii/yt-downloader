package main

import (
	"encoding/json"
	"errors"
	"os"
)

// TODO move somewhere to temp dir
var defaultSettingsFileLocation = "./settings.json"

type AppSettings struct {
	SaveDirectory string `json:"saveDirectory"`
}

// loads existing app settings, creates settings file if settings file not found
func LaodAppSettings() (*AppSettings, error) {
	_, err := os.Stat(defaultSettingsFileLocation)
	if errors.Is(err, os.ErrNotExist) {
		settings := &AppSettings{
			SaveDirectory: "./yt-downloader",
		}
		data, err := json.Marshal(settings)
		if err != nil {
			return nil, err
		}
		file, err := os.Create(defaultSettingsFileLocation)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			return nil, err
		}
		return settings, nil
	}
	data, err := os.ReadFile(defaultSettingsFileLocation)
	if err != nil {
		return nil, err
	}
	settings := &AppSettings{}
	err = json.Unmarshal(data, settings)
	if err != nil {
		return nil, err
	}
	return settings, nil
}
