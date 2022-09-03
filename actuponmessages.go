package main

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var UserIdAlreadyDealtWith = make(map[int]map[int]bool)

func IsMessageGDPR(Message MessageStruct) bool {
	return Message.IsSystemMessage && Message.Subject == "[Important] Right to Erasure - Action Requested"
}

func GetGDPRInfoFromMessage(Message MessageStruct) ([]int, []int) {
	// UserIds := make([]int, 0)
	// PlaceIds := make([]int, 0)
	Body := Message.Body

	UserIdsString := strings.Split(strings.Split(Body, "Please delete this User ID")[0], "following User ID(s):")[1]
	PlaceIdsString := strings.Split(strings.Split(Body, "This is an obligation under data protection laws.")[0], "from the following Game(s):")[1]

	re := regexp.MustCompile("[0-9]+")

	UserIds := re.FindAllString(UserIdsString, -1)
	PlaceIds := re.FindAllString(PlaceIdsString, -1)

	FilteredUserIds := make([]int, 0)
	FilteredPlaceIds := make([]int, 0)

	var PlaceIdsMap = make(map[string]bool)

	for _, PlaceId := range PlaceIds {
		if PlaceId != "160" && !PlaceIdsMap[PlaceId] {
			PlaceIdInt, _ := strconv.Atoi(PlaceId)
			PlaceIdsMap[PlaceId] = true
			FilteredPlaceIds = append(FilteredPlaceIds, PlaceIdInt)
		}
	}

	for _, UserId := range UserIds {
		if UserId != "160" {
			UserIdInt, _ := strconv.Atoi(UserId)
			FilteredUserIds = append(FilteredUserIds, UserIdInt)
		}
	}

	bodyByteArrayUserIds, _ := json.Marshal(FilteredUserIds)
	bodyByteArrayPlaceIds, _ := json.Marshal(FilteredPlaceIds)

	println("UserIds")
	println(UserIdsString)
	println(string(bodyByteArrayUserIds))

	println("PlaceIds")
	println(PlaceIdsString)
	println(string(bodyByteArrayPlaceIds))

	return FilteredPlaceIds, FilteredUserIds
}

func ActUponMessages(Messages []MessageStruct) {
	for _, Message := range Messages {
		if !IsMessageGDPR(Message) {
			continue
		}

		PlaceIds, UserIds := GetGDPRInfoFromMessage(Message)
		GetGameIdFromPlaceIds(PlaceIds)

		for _, PlaceId := range PlaceIds {
			PlaceIdStr := strconv.Itoa(PlaceId)

			if DataKeys[PlaceIdStr] == nil {
				println(PlaceIdStr + " has no data key! info")
				continue
			}

			Success, GameId := GetGameIdFromPlaceId(PlaceId)

			if !Success {
				println("FAILED TO GET GAMEID FOR " + PlaceIdStr)
				continue
			}

			DataStoreNames := DataKeys[PlaceIdStr]

			for _, UserId := range UserIds {
				if UserIdAlreadyDealtWith[PlaceId][UserId] {
					continue
				}

				UserIdStr := strconv.Itoa(UserId)

				CompleteSuccess := true

				for DataStoreName, Keys := range DataStoreNames {
					PlaceDataStore := DataStore{PlaceId: PlaceId, GameId: GameId, Name: DataStoreName}

					for _, Key := range Keys {
						for {
							Success, Response := PlaceDataStore.RemoveAsync(strings.Replace(Key, `%USERID`, UserIdStr, -1))
							StatusCode := Response.StatusCode

							if !Success {
								if StatusCode == 404 {
									break
								} else if StatusCode == 429 {
									time.Sleep(time.Second * 10)
									continue
								}

								CompleteSuccess = false
								println("FAILED TO DELETE KEY FOR " + UserIdStr)
								println(StatusCode)
							}

							break
						}
					}
				}

				if CompleteSuccess {
					if UserIdAlreadyDealtWith[PlaceId] == nil {
						UserIdAlreadyDealtWith[PlaceId] = make(map[int]bool)
					}

					UserIdAlreadyDealtWith[PlaceId][UserId] = true
				}
			}
		}
		break
	}
}
