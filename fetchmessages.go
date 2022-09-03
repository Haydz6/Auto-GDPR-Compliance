package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type DeleteMessagesBodyReqStruct struct {
	MessageIds []int `json:"messageIds"`
}

type MessageSenderStruct struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type MessageRecipientStruct struct {
	Id          int    `json:""`
	Name        string `json:""`
	DisplayName string `json:""`
}

type MessageStruct struct {
	Id                     int                    `json:"id"`
	Sender                 MessageSenderStruct    `json:"sender"`
	Recipient              MessageRecipientStruct `json:"recipient"`
	Subject                string                 `json:"subject"`
	Body                   string                 `json:"body"`
	Created                string                 `json:"created"`
	Updated                string                 `json:"updated"`
	IsRead                 bool                   `json:"isRead"`
	IsSystemMessage        bool                   `json:"isSystemMessage"`
	IsReportAbuseDisplayed bool                   `json:"isReportAbuseDisplayed"`
}

type MessagesStruct struct {
	Collection          []MessageStruct `json:"collection"`
	TotalCollectionSize int             `json:"totalCollectionSize"`
	TotalPages          int             `json:"totalPages"`
	PageNumber          int             `json:"pageNumber"`
}

func FetchMessages(PageNum int) (bool, *http.Response, bool, int, []MessageStruct) {
	Success, Response := RobloxRequest(fmt.Sprintf("https://privatemessages.roblox.com/v1/messages?pageNumber=%d&pageSize=20&messageTab=Inbox", PageNum), "GET", nil, "")

	if !Success {
		if Response.StatusCode == 401 {
			println("Your ROBLOSECURITY is not valid!")
			time.Sleep(time.Second * 3)
			os.Exit(1)
		}

		println("Failed to fetch messages!")
		println(Response.StatusCode)
		return false, Response, false, 0, nil
	}

	var Body MessagesStruct
	json.NewDecoder(Response.Body).Decode(&Body)

	return true, Response, Body.PageNumber+1 >= Body.TotalPages, Body.PageNumber, Body.Collection
}

func DeleteMessages(MessageIds []int) (bool, *http.Response) {
	bodyByteArray, err := json.Marshal(DeleteMessagesBodyReqStruct{MessageIds: MessageIds})

	if err != nil {
		println(err.Error())
		return false, nil
	}

	Success, Response := RobloxRequest("https://privatemessages.roblox.com/v1/messages/archive", "POST", nil, string(bodyByteArray))

	if !Success {
		println("Failed to delete messages!")
		println(Response.StatusCode)
		return false, Response
	}

	return true, Response
}

func ReadMessages(MessageIds []int) (bool, *http.Response) {
	bodyByteArray, err := json.Marshal(DeleteMessagesBodyReqStruct{MessageIds: MessageIds})

	if err != nil {
		println(err.Error())
		return false, nil
	}

	Success, Response := RobloxRequest("https://privatemessages.roblox.com/v1/messages/mark-read", "POST", nil, string(bodyByteArray))

	if !Success {
		println("Failed to mark read messages!")
		println(Response.StatusCode)
		return false, Response
	}

	return true, Response
}
