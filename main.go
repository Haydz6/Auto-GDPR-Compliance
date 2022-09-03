package main

import (
	"os"
	"path"
	"strings"
	"time"
)

var ROBLOSECURITY = ""

func GetROBLOSecurity() {
	String, err := os.ReadFile(path.Join(".", "ROBLOSECURITY.txt"))

	if err != nil {
		panic(err.Error())
	}

	ROBLOSECURITY = strings.ReplaceAll(string(String), "Put your ROBLOSECURITY below!!!", "")
}

func main() {
	GetROBLOSecurity()
	FetchDataKeys()

	if ROBLOSECURITY == "" {
		println("You did not put your ROBLOSECURITY in!")
		println("Exiting in 3 seconds")
		time.Sleep(time.Second * 3)
		return
	}

	PreviousPageNumber := 1
	for {
		Success, Response, IsEnd, PageNumber, Messages := FetchMessages(PreviousPageNumber)
		PreviousPageNumber = PageNumber + 1

		if !Success {
			println("Failed to fetch message!")
			println(Response.StatusCode)
			break
		}

		ActUponMessages(Messages)

		if IsEnd {
			println("Went through all the messages!")
			break
		}

		break
	}

	println("Done!")
	time.Sleep(time.Second * 3)
}
