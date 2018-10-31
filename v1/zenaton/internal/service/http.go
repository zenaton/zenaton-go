package service

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var client = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
}

// Get sends a GET request to the specified url
func Get(url string) (*http.Response, error) {
	return client.Get(url)
}

// Post sends a json POST http request to the specified url with the specified body
func Post(url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
}

// Put sends a json PUT http request to the specified url with the specified body
func Put(url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}
