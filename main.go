package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var getaroundBaseURL = "https://api.getaround.com/"

type GetaroundClient struct {
	token string

	baseURL    string
	httpClient *http.Client
}

func main() {
	fmt.Println("hello world")

	token := os.Getenv("GETAROUND_ACCESS_TOKEN")
	if token == "" {
		panic(errors.New("you must provide an access token"))
	}

	client := GetaroundClient{
		token:      token,
		baseURL:    getaroundBaseURL,
		httpClient: http.DefaultClient,
	}

	data, err := client.getUser()
	if err != nil {
		panic(err)
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(pretty))
}

func (c *GetaroundClient) do(req *http.Request) (map[string]any, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var generic map[string]any
	err = json.Unmarshal(data, &generic)
	if err != nil {
		return nil, err
	}

	return generic, nil
}

func (c *GetaroundClient) getUser() (map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"users/me", nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}
