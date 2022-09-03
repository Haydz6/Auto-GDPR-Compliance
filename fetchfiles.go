package main

import (
	"encoding/json"
	"os"
	"path"
	"time"
)

type SettingsStruct struct {
	DeleteGDPRMessagesAfterFulfilled     bool
	MarkGDPRMessagesAsReadAfterFulfilled bool
}

var DataKeys map[string]map[string]map[string]map[string][]string
var Settings SettingsStruct
var HandledMessageIds map[int]bool

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

func FetchHandledIds() {
	Bytes, err := os.ReadFile(path.Join(".", "HandledIds.json"))

	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(Bytes, &HandledMessageIds)
}

func SaveHandledIds() {
	Bytes, err := json.Marshal(HandledMessageIds)

	if err != nil {
		println(err.Error())
		return
	}

	WriteErr := os.WriteFile(path.Join(".", "HandledIds.json"), Bytes, 0644)

	if WriteErr != nil {
		println(WriteErr.Error())
	}
}

func LoopSaveHandleIds() {
	go func() {
		for {
			time.Sleep(time.Second * 10)
			SaveHandledIds()
		}
	}()
}
