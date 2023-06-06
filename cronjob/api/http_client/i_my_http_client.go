package http_client

import (
	"net/http"
)

type IMyHttpClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, jsonBody []byte) error
	Put(url string, jsonBody []byte) (*http.Response, error)
	Delete(url string) error
}
