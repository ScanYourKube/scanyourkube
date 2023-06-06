package api_keel_webhook_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/arbitrary"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	api_keel_webhook "github.com/scanyourkube/cronjob/api/keel"
	update_keel "github.com/scanyourkube/cronjob/dto/api/update"
)

func TestUpdateImage(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1024)
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	arbitraries := arbitrary.DefaultArbitraries()
	arbitraries.RegisterGen(gen.Struct(reflect.TypeOf(&update_keel.KeelWebhookDto{}), map[string]gopter.Gen{
		"Name": gen.AlphaString().SuchThat(func(s string) bool {
			return s != ""
		}),
		"Tag": gen.AlphaString(),
	}).Map(func(dto update_keel.KeelWebhookDto) update_keel.KeelWebhookDto {
		return dto
	}))

	properties.Property("returns valid JSON when provided with a valid keelWebhookDto object", prop.ForAll(
		func(keelWebhookDto update_keel.KeelWebhookDto) bool {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": "ok"}`))
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			a := api_keel_webhook.NewKeelWebHookApi(server.Client(), server.URL)
			err := a.UpdateImage(keelWebhookDto)

			return err == nil
		},
		arbitraries.GenForType(reflect.TypeOf(update_keel.KeelWebhookDto{})),
	))

	properties.Property("returns an error when provided with an invalid keelWebhookDto object", prop.ForAll(
		func(keelWebhookDto update_keel.KeelWebhookDto) bool {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status": "error"}`))
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			a := api_keel_webhook.NewKeelWebHookApi(server.Client(), server.URL)
			err := a.UpdateImage(keelWebhookDto)

			return err != nil
		},
		arbitraries.GenForType(reflect.TypeOf(update_keel.KeelWebhookDto{})),
	))

	properties.Property("returns a status code of either 200 or 409 when provided with a valid keelWebhookDto object", prop.ForAll(
		func(keelWebhookDto update_keel.KeelWebhookDto) bool {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": "ok"}`))
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			a := api_keel_webhook.NewKeelWebHookApi(server.Client(), server.URL)
			err := a.UpdateImage(keelWebhookDto)

			return err == nil
		},
		arbitraries.GenForType(reflect.TypeOf(update_keel.KeelWebhookDto{})),
	))

	properties.Property("returns an error when the client is unable to make a POST request to the specified URL", prop.ForAll(
		func(keelWebhookDto update_keel.KeelWebhookDto) bool {

			a := api_keel_webhook.NewKeelWebHookApi(&http.Client{}, "http://invalid-url")
			err := a.UpdateImage(keelWebhookDto)

			return err != nil
		},
		arbitraries.GenForType(reflect.TypeOf(update_keel.KeelWebhookDto{})),
	))

	properties.TestingRun(t)
}
