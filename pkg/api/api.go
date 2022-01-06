package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/miladibra10/vjson"
	log "github.com/sirupsen/logrus"
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

func (c *ApiClient) Submit(method string, path string, data string, schema string) error {

	err := c.validate(data, schema)
	if err != nil {
		fmt.Println("Error validating json: ", err, schema)
		return err
	}
	fmt.Printf("path: %+v\n", path)
	//fmt.Printf("sendRequest: %+v\n", data)
	//return nil

	//TODO make TLS verification configurable
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(method, c.apiHost+path, strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.apiToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, clienterr := client.Do(req)
	if clienterr != nil {
		return clienterr
	}
	defer resp.Body.Close()

	if resp.StatusCode > 201 {
		log.Debugf("response: %+v\n", resp)
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(b))
		var result struct {
			Message   string `json:"message"`
			Exception string `json:"exception"`
			File      string `json:"file"`
			Line      int    `json:"line"`
			Trace     []struct {
				File     string `json:"file"`
				Line     int    `json:"line"`
				Function string `json:"function"`
				Class    string `json:"class"`
				Type     string `json:"type"`
			} `json:"trace"`
		}

		json.Unmarshal(b, &result)
		log.Errorf("response: %+s\n\n", result.Message)
		log.Panic("Error submitting to API: ", resp.Status)
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
