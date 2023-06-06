package http_client

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyHttpClientGet(t *testing.T) {
	// Create a test HTTP client with a mock Do function
	mockClient := &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify that the request is a GET request
				assert.Equal(t, http.MethodGet, req.Method)

				// Create a mock response
				response := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("OK")),
				}

				return response, nil
			},
		},
	}

	// Create an instance of MyHttpClient using the mock client
	client := NewMyHttpClient(mockClient)

	// Perform the GET request
	response, err := client.Get("http://example.com")

	// Verify the response and error
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "OK", string(body))
}

func TestMyHttpClientPost(t *testing.T) {
	// Create a test HTTP client with a mock Do function
	mockClient := &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify that the request is a POST request
				assert.Equal(t, http.MethodPost, req.Method)

				// Verify the request body
				body, _ := io.ReadAll(req.Body)
				assert.Equal(t, []byte(`{"key":"value"}`), body)

				// Create a mock response
				response := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("OK")),
				}

				return response, nil
			},
		},
	}

	// Create an instance of MyHttpClient using the mock client
	client := NewMyHttpClient(mockClient)

	// Perform the POST request
	err := client.Post("http://example.com", []byte(`{"key":"value"}`))

	// Verify the error
	assert.NoError(t, err)
}

func TestMyHttpClientPut(t *testing.T) {
	// Create a test HTTP client with a mock Do function
	mockClient := &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify that the request is a PUT request
				assert.Equal(t, http.MethodPut, req.Method)

				// Verify the request body
				body, _ := io.ReadAll(req.Body)
				assert.Equal(t, []byte(`{"key":"value"}`), body)

				// Create a mock response
				response := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("OK")),
				}

				return response, nil
			},
		},
	}

	// Create an instance of MyHttpClient using the mock client
	client := NewMyHttpClient(mockClient)

	// Perform the PUT request
	response, err := client.Put("http://example.com", []byte(`{"key":"value"}`))

	// Verify the response and error
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "OK", string(body))
}

func TestMyHttpClientDelete(t *testing.T) {
	// Create a test HTTP client with a mock Do function
	mockClient := &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify that the request is a DELETE request
				assert.Equal(t, http.MethodDelete, req.Method)

				// Create a mock response
				response := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("OK")),
				}

				return response, nil
			},
		},
	}

	// Create an instance of MyHttpClient using the mock client
	client := NewMyHttpClient(mockClient)

	// Perform the DELETE request
	err := client.Delete("http://example.com")

	// Verify the error
	assert.NoError(t, err)
}

// mockTransport is a custom RoundTripper that allows mocking HTTP requests and responses
type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTripFunc(req)
}
