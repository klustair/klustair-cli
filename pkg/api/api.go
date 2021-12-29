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

func (c *ApiClient) validate(json []byte, schemapath string) error {

	schema, err := vjson.ReadFromFile("./pkg/api/schema/" + schemapath + ".json")
	if err != nil {
		fmt.Println("Error reading schema: ", err)
		return err
	}

	err = schema.ValidateBytes(json)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) SendReport(report []byte) error {

	err := c.validate(report, "report")
	if err != nil {
		return err
	}
	err = c.sendRequest("POST", "/api/v1/pac/report/create", string(report))
	return err
}
