package main

import (
	"encoding/json"
	"os"
	"path"
)

type SettingsStruct struct {
	DeleteGDPRMessagesAfterFulfilled     bool
	MarkGDPRMessagesAsReadAfterFulfilled bool
}

var DataKeys map[string]map[string]map[string][]string
var Settings SettingsStruct

func FetchDataKeys() {
	Bytes, err := os.ReadFile(path.Join(".", "DataKeys.json"))

	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(Bytes, &DataKeys)
}

func FetchSettings() {
	Bytes, err := os.ReadFile(path.Join(".", "Settings.json"))

	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(Bytes, &Settings)
}
