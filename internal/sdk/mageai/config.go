package mageai

import "net/http"

type ClientConfig struct {
	ApiKey     string
	Host       string
	HTTPClient *http.Client
}
