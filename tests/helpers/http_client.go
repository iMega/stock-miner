package helpers

import (
	"net/http"
	"net/http/cookiejar"
)

var httpClient *http.Client

func GetHTTPClient() *http.Client {
	if httpClient != nil {
		return httpClient
	}

	jar, _ := cookiejar.New(nil)
	httpClient = &http.Client{Jar: jar}
	httpClient.Get("http://app/google/callback")

	return httpClient
}
