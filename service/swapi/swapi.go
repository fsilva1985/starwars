package swapi

import (
	"net/http"
	"time"
)

type Client interface {
	Do(request *http.Request) (*http.Response, error)
	Url() string
}

func NewClient() Client {
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	return &httpClientImpl{
		url:    "https://swapi.dev/api/",
		client: &http.Client{Transport: tr},
	}
}

type httpClientImpl struct {
	url    string
	client *http.Client
}

type body struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  interface{} `json:"results"`
}

func (h *httpClientImpl) Url() string {
	return h.url
}

func (h *httpClientImpl) Do(request *http.Request) (*http.Response, error) {
	response, err := h.client.Do(request)

	return response, err
}
