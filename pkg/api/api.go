package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/miladibra10/vjson"
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

func (c *ApiClient) Submit(method string, path string, json string, schema string) error {

	err := c.validate(json, schema)
	if err != nil {
		fmt.Println("Error validating json: ", err, schema)
		return err
	}
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

func (c *ApiClient) validate(json string, schema string) error {

	sma, err := vjson.ReadFromFile("./pkg/api/schema/" + schema + ".json")
	if err != nil {
		fmt.Println("Error reading schema: ", err)
		return err
	}

	err = sma.ValidateString(json)
	if err != nil {
		return err
	}
	return nil
}
