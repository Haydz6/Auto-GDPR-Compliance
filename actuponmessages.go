package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var UserIdAlreadyDealtWith = make(map[int]map[int]bool)
var NoValidDataKeysEntry = make(map[int]bool)

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

	// bodyByteArrayUserIds, _ := json.Marshal(FilteredUserIds)
	// bodyByteArrayPlaceIds, _ := json.Marshal(FilteredPlaceIds)

	// println("UserIds")
	// println(UserIdsString)
	// println(string(bodyByteArrayUserIds))

	// println("PlaceIds")
	// println(PlaceIdsString)
	// println(string(bodyByteArrayPlaceIds))

	return FilteredPlaceIds, FilteredUserIds
}

func ActUponMessages(Messages []MessageStruct) {
	SuccessfulMessageIds := make([]int, 0)

	for _, Message := range Messages {
		if !IsMessageGDPR(Message) {
			continue
		}

		PlaceIds, UserIds := GetGDPRInfoFromMessage(Message)
		GetGameIdFromPlaceIds(PlaceIds)

		AllSuccessful := true

		for _, PlaceId := range PlaceIds {
			PlaceIdStr := strconv.Itoa(PlaceId)

			if DataKeys[PlaceIdStr] == nil {

				if !NoValidDataKeysEntry[PlaceId] {
					NoValidDataKeysEntry[PlaceId] = true
					println(PlaceIdStr + " does not have a data key entry!")
				}

				AllSuccessful = false
				continue
			}

			Success, GameId := GetGameIdFromPlaceId(PlaceId)

			if !Success {
				AllSuccessful = false
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

				for DataStoreName, Scopes := range DataStoreNames {
					for Scope, Keys := range Scopes {
						PlaceDataStore := DataStore{PlaceId: PlaceId, GameId: GameId, Name: DataStoreName, Scope: Scope}

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
							}
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

		if AllSuccessful && (Settings.DeleteGDPRMessagesAfterFulfilled || Settings.MarkGDPRMessagesAsReadAfterFulfilled) {
			SuccessfulMessageIds = append(SuccessfulMessageIds, Message.Id)
		}
	}

	if len(SuccessfulMessageIds) > 0 {
		if Settings.DeleteGDPRMessagesAfterFulfilled {
			DeleteMessages(SuccessfulMessageIds)
		} else if Settings.MarkGDPRMessagesAsReadAfterFulfilled {
			ReadMessages(SuccessfulMessageIds)
		}
	}
}
