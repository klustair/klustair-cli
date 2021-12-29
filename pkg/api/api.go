package api

import (
	"fmt"
	"net/http"
	"strings"
)

type ApiClient struct {
	apiHost  string
	apiToken string
}

func NewApiClient(apiHost string, apiToken string) *ApiClient {
	return &ApiClient{
		apiHost:  apiHost,
		apiToken: apiToken,
	}
}

func (c *ApiClient) sendRequest(method string, path string, json string) error {

	fmt.Printf("sendRequest: %+v\n", json)
	return nil

	req, err := http.NewRequest(method, c.apiHost+path, strings.NewReader(json))
	req.Header.Add("Authorization", "Bearer "+c.apiToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) SendReport(report []byte) error {
	err := c.sendRequest("POST", "/api/v1/reports", string(report))
	return err
}
