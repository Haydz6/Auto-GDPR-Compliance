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
	Type    string
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
	if Store.Type == "" {
		Store.Type = "DataStore"
	}

	if Store.Type == "DataStore" {
		return RobloxRequest(fmt.Sprintf("https://gamepersistence.roblox.com/v2/persistence/%d/datastores/objects/object?datastore=%s&objectKey=%s%s%s", Store.GameId, Store.Name, Store.Scope, "%2F", key), "GET", map[string]string{"Roblox-Place-Id": strconv.Itoa(Store.PlaceId)}, "")
	} else if Store.Type == "OrderedDataStore" {
		return RobloxRequest(fmt.Sprintf("https://gamepersistence.roblox.com/v1/persistence/sorted?scope=%s&key=%s&target=%s", Store.Scope, Store.Name, key), "GET", map[string]string{"Roblox-Place-Id": strconv.Itoa(Store.PlaceId)}, "")
	}

	return false, nil
}

func (Store *DataStore) SetAsync(key string, body interface{}) (bool, *http.Response) {
	if Store.Scope == "" {
		Store.Scope = "global"
	}
	if Store.Type == "" {
		Store.Type = "DataStore"
	}

	bodyByteArray, _ := json.Marshal(body)

	if Store.Type == "DataStore" {
		digest := md5.New()
		digest.Write(bodyByteArray)
		bodyHash := digest.Sum(nil)
		bodyHashBase64 := base64.StdEncoding.EncodeToString(bodyHash)

		return RobloxRequest(fmt.Sprintf("https://gamepersistence.roblox.com/v2/persistence/%d/datastores/objects/object?datastore=%s&objectKey=%s%s%s", Store.GameId, Store.Name, Store.Scope, "%2F", key), "POST", map[string]string{"Roblox-Place-Id": strconv.Itoa(Store.PlaceId), "Content-Type": "application/octet-stream", "Content-MD5": bodyHashBase64}, string(bodyByteArray))
	} else if Store.Type == "OrderedDataStore" {
		return RobloxRequest(fmt.Sprintf("https://gamepersistence.roblox.com/v1/persistence/sorted?scope=%s&key=%s&target=%s", Store.Scope, Store.Name, key), "POST", map[string]string{"Roblox-Place-Id": strconv.Itoa(Store.PlaceId), "Content-Type": "application/octet-stream"}, string(bodyByteArray))
	}

	return false, nil
}

func (Store *DataStore) RemoveAsync(key string) (bool, *http.Response) {
	if Store.Scope == "" {
		Store.Scope = "global"
	}
	if Store.Type == "" {
		Store.Type = "DataStore"
	}

	if Store.Type == "DataStore" {
		return RobloxRequest(fmt.Sprintf("https://gamepersistence.roblox.com/v2/persistence/%d/datastores/objects/object?datastore=%s&objectKey=%s%s%s", Store.GameId, Store.Name, Store.Scope, "%2F", key), "DELETE", map[string]string{"Roblox-Place-Id": strconv.Itoa(Store.PlaceId)}, "")
	} else if Store.Type == "OrderedDataStore" {
		return RobloxRequest(fmt.Sprintf("https://gamepersistence.roblox.com/v1/persistence/sorted/remove?scope=%s&key=%s&target=%s", Store.Scope, Store.Name, key), "POST", map[string]string{"Roblox-Place-Id": strconv.Itoa(Store.PlaceId)}, "")
	}

	return false, nil
}

func GetDataStore(GameId int, PlaceId int, Type string, Name string, Scope string) DataStore {
	return DataStore{GameId, PlaceId, Type, Name, Scope}
}
