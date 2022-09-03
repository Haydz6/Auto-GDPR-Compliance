package main

import (
	"encoding/json"
	"os"
	"path"
)

var DataKeys map[string]map[string]map[string][]string

func FetchDataKeys() {
	Bytes, err := os.ReadFile(path.Join(".", "DataKeys.json"))

	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(Bytes, &DataKeys)
}
