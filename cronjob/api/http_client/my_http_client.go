package http_client

import (
	"bytes"
	"net/http"
)

type MyHttpClient struct {
	client *http.Client
}

func NewMyHttpClient(client *http.Client) MyHttpClient {
	return MyHttpClient{client: client}
}

func (h MyHttpClient) do(req *http.Request) (*http.Response, error) {
	return h.client.Do(req)
}

func (h MyHttpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := h.do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (h MyHttpClient) Post(url string, jsonBody []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = h.do(req)

	if err != nil {
		return err
	}

	return nil
}

func (h MyHttpClient) Put(url string, jsonBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := h.do(req)

	if err != nil {
		return response, err
	}

	return response, nil
}

func (h MyHttpClient) Delete(url string) error {
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	_, err = h.do(req)

	if err != nil {
		return err
	}

	return nil
}
