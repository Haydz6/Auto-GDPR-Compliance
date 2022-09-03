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

	TextSplit := strings.Split(string(String), "\n")

	if len(TextSplit) > 1 {
		ROBLOSECURITY = TextSplit[1]
	} else {
		ROBLOSECURITY = TextSplit[0]
	}
}

func main() {
	GetROBLOSecurity()
	FetchDataKeys()
	FetchSettings()
	FetchHandledIds()
	LoopSaveHandleIds()

	if ROBLOSECURITY == "" {
		println("You did not put your ROBLOSECURITY in!")
		println("Exiting in 3 seconds")
		time.Sleep(time.Second * 3)
		return
	}

	println("Running")

	PreviousPageNumber := 0
	for {
		Completed := false

		for {
			Success, Response, IsEnd, PageNumber, Messages := FetchMessages(PreviousPageNumber)
			PreviousPageNumber = PageNumber + 1

			if !Success {
				StatusCode := Response.StatusCode

				if StatusCode == 429 {
					time.Sleep(time.Second * 10)
					break
				}

				println("Failed to fetch message!")
				println(StatusCode)
				break
			}

			ActUponMessages(Messages)

			Completed = IsEnd

			if IsEnd {
				println("Went through all the messages!")
				break
			}
		}

		if Completed {
			break
		}
	}

	SaveHandledIds()

	println("Done!")
	time.Sleep(time.Second * 3)
}
