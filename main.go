package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var getaroundBaseURL = "https://api.getaround.com/"
var debug = true

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

	user, err := client.getUser()
	if err != nil {
		panic(err)
	}

	cars, err := client.listCars(user["id"].(string))
	if err != nil {
		panic(err)
	}

	if cars["count"].(float64) == 0 {
		log.Println("No cars found, exiting")
		return
	}
	carID := cars["items"].([]any)[0].(map[string]any)["id"].(string)

	data, err := client.tripsByCar(carID)
	if err != nil {
		panic(err)
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(pretty))
}

func (c *GetaroundClient) getUser() (map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"users/me", nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *GetaroundClient) listCars(ownerID string) (map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"cars", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("owner_id", ownerID)
	req.URL.RawQuery = q.Encode()

	return c.do(req)
}

func (c *GetaroundClient) tripsByCar(carID string) (map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"trips", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("car_id", carID)
	q.Add("sort", "starts_at")
	q.Add("summary_renter_trip_status", "ENDED")
	q.Add("fields", "items{id,reservation,status}")
	req.URL.RawQuery = q.Encode()

	return c.do(req)
}

func (c *GetaroundClient) do(req *http.Request) (map[string]any, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	log.Printf(fmt.Sprintf("%s %s%s", req.Method, req.URL.Host, req.URL.Path))
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
