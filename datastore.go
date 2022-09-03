package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DataStore struct {
	GameId  int
	PlaceId int
	Name    string
	Scope   string
}

type DataStoreErrResponse struct {
	errors []struct {
		code      int
		message   string
		retryable bool
	}
}

func (Store *DataStore) GetAsync(key string) (bool, *http.Response) {
	if Store.Scope == "" {
		Store.Scope = "global"
	}

	Success, Response := RobloxRequest(fmt.Sprintf("https://gamepersistence.roblox.com/v2/persistence/%d/datastores/objects/object?datastore=%s&objectKey=%s%s%s", Store.GameId, Store.Name, Store.Scope, "%2F", key), "GET", map[string]string{"Roblox-Place-Id": strconv.Itoa(Store.PlaceId)}, "")

	return Success, Response
}

func (Store *DataStore) SetAsync(key string, body interface{}) (bool, *http.Response) {
	if Store.Scope == "" {
		Store.Scope = "global"
	}

	bodyByteArray, _ := json.Marshal(body)
	digest := md5.New()
	digest.Write(bodyByteArray)
	bodyHash := digest.Sum(nil)
	bodyHashBase64 := base64.StdEncoding.EncodeToString(bodyHash)

	Success, Response := RobloxRequest(fmt.Sprintf("https://gamepersistence.roblox.com/v2/persistence/%d/datastores/objects/object?datastore=%s&objectKey=%s%s%s", Store.GameId, Store.Name, Store.Scope, "%2F", key), "POST", map[string]string{"Roblox-Place-Id": strconv.Itoa(Store.PlaceId), "Content-Type": "application/octet-stream", "Content-MD5": bodyHashBase64}, string(bodyByteArray))

	return Success, Response
}

func (Store *DataStore) RemoveAsync(key string) (bool, *http.Response) {
	if Store.Scope == "" {
		Store.Scope = "global"
	}

	Success, Response := RobloxRequest(fmt.Sprintf("https://gamepersistence.roblox.com/v2/persistence/%d/datastores/objects/object?datastore=%s&objectKey=%s%s%s", Store.GameId, Store.Name, Store.Scope, "%2F", key), "DELETE", map[string]string{"Roblox-Place-Id": strconv.Itoa(Store.PlaceId)}, "")

	return Success, Response
}

func GetDataStore(GameId int, PlaceId int, Name string, Scope string) DataStore {
	return DataStore{GameId, PlaceId, Name, Scope}
}
