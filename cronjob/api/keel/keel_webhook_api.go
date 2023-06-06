package api_keel_webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	update_keel "github.com/scanyourkube/cronjob/dto/api/update"
)

type IKeelWebhookApi interface {
	UpdateImage(keelWebhookDto update_keel.KeelWebhookDto) error
}

type KeelWebHookApi struct {
	client  *http.Client
	baseUrl string
}

func NewKeelWebHookApi(client *http.Client, baseUrl string) IKeelWebhookApi {
	return KeelWebHookApi{client: client, baseUrl: baseUrl}
}

func (a KeelWebHookApi) UpdateImage(keelWebhookDto update_keel.KeelWebhookDto) error {
	json, _ := json.Marshal(keelWebhookDto)

	response, err := a.client.Post(a.baseUrl+"/v1/webhooks/native", "application/json; charset=utf-8", bytes.NewBuffer(json))

	if err != nil {
		return err
	}

	if response.StatusCode != 200 && response.StatusCode != 409 {
		return fmt.Errorf("returned response code didn't show a successful state status code: %d", response.StatusCode)
	}

	return nil
}
