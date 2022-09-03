package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
		println("Failed to fetch messages!")
		println(Response.StatusCode)
		return false, Response, false, 0, nil
	}

	var Body MessagesStruct
	json.NewDecoder(Response.Body).Decode(&Body)

	return true, Response, Body.PageNumber+1 >= Body.TotalPages, Body.PageNumber, Body.Collection
}
