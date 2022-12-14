package main

import (
	"encoding/json"
	"io"
	"strconv"
)

var GameIdLookup = make(map[int]int)

type GameDetail struct {
	PlaceId               int    `json:"placeId"`
	Name                  string `json:"name"`
	SourceName            string `json:"sourceName"`
	SourceDescription     string `json:"sourceDescription"`
	Url                   string `json:"url"`
	Builder               string `json:"builder"`
	BuilderId             int    `json:"builderId"`
	HasVerifiedBadge      bool   `json:"hasVerifiedBadge"`
	IsPlayable            bool   `json:"isPlayable"`
	ReasonProhibited      string `json:"reasonProhibited"`
	UniverseId            int    `json:"universeId"`
	UniverseIdRootPlaceId int    `json:"universeIdRootPlaceId"`
	Price                 int    `json:"price"`
	ImageToken            string `json:"imageToken"`
}

type GameDetailsReq struct {
	UniverseIds []int `json:"universeIds"`
}

func GetGameIdFromPlaceId(PlaceId int) (bool, int) {
	CachedGameId := GameIdLookup[PlaceId]

	if CachedGameId != 0 {
		return true, CachedGameId
	}

	Success, Response := RobloxRequest("https://games.roblox.com/v1/games/multiget-place-details?placeIds="+strconv.Itoa(PlaceId), "GET", nil, "")

	if !Success {
		println("Failed to fetch gameid!")
		println(Response.StatusCode)
		println(PlaceId)
		println("https://games.roblox.com/v1/games/multiget-place-details?placeIds=" + strconv.Itoa(PlaceId))

		bodyBytes, err := io.ReadAll(Response.Body)
		if err == nil {
			println(string(bodyBytes))
		}

		return false, 0
	}

	var Data []GameDetail
	json.NewDecoder(Response.Body).Decode(&Data)

	Game := Data[0]

	GameId := Game.UniverseId

	GameIdLookup[Game.UniverseIdRootPlaceId] = GameId

	return true, GameId
}

func GetGameIdFromPlaceIds(PlaceIds []int) (bool, map[int]int) {
	GameIds := make(map[int]int)
	//GameIdsReq := ""

	for _, PlaceId := range PlaceIds {
		GameId := GameIdLookup[PlaceId]
		if GameId != 0 {
			GameIds[GameId] = PlaceId
		} else {
			// if GameIdsReq == "" {
			// 	GameIdsReq = string(GameId)
			// } else {
			// 	GameIdsReq = GameIdsReq + "%" + "2C" + string(GameId)
			// }
			Success, GameId := GetGameIdFromPlaceId(PlaceId)

			if !Success {
				println("Failed to get placeid from gameid!")
				continue
			}

			GameIds[GameId] = PlaceId
			GameIdLookup[PlaceId] = GameId
		}
	}

	// Success, Response := RobloxRequest("https://games.roblox.com/v1/games?universeIds="+GameIdsReq, "GET", nil, "")

	// if !Success {
	// 	println("Failed to fetch gameids!")
	// 	println(Response.StatusCode)
	// 	return false, nil
	// }

	// var Data []GameDetail
	// json.NewDecoder(Response.Body).Decode(&Data)

	// for _, Game := range Data {
	// 	GameId := Game.UniverseId
	// 	PlaceId := Game.UniverseIdRootPlaceId

	// 	GameIds[GameId] = PlaceId
	// 	GameIdLookup[PlaceId] = GameId
	// }

	return true, GameIds
}
