package main

import (
	"bytes"
	"fmt"
	"net/http"
)

var CSRFToken = ""

func RobloxRequest(URL string, Method string, Headers map[string]string, Body string) (bool, *http.Response) {
	client := &http.Client{}
	request, reqerr := http.NewRequest(Method, URL, bytes.NewBuffer([]byte(Body)))

	if reqerr != nil {
		panic(reqerr.Error())
	}

	for k, v := range Headers {
		request.Header.Set(k, v)
	}

	request.Header.Set("Referer", "https://www.roblox.com")
	request.Header.Set("Origin", "https://www.roblox.com")
	request.Header.Set("Cookie", fmt.Sprintf(".ROBLOSECURITY=%s", ROBLOSECURITY))
	request.Header.Set("X-CSRF-Token", CSRFToken)

	if Headers["Content-Type"] == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := client.Do(request)

	if err != nil {
		println(ROBLOSECURITY)
		println(err.Error())
	}

	if response.StatusCode == 403 {
		ICSRFToken := response.Header.Get("x-csrf-token")
		if ICSRFToken != "" {
			CSRFToken = ICSRFToken
			return RobloxRequest(URL, Method, Headers, Body)
		}
	}

	return err == nil && response.StatusCode < 400, response
}
