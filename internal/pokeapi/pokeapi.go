package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const BaseURL = "https://pokeapi.co/api/v2/location-area"

var client = &http.Client{
	Timeout: time.Second * 23,
}

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// get locations area struct from url
func GetLocations(url string) (LocationAreasResponse, error) {
	if url == "" {
		return LocationAreasResponse{}, fmt.Errorf("url is empty")
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer resp.Body.Close()
	var locResp LocationAreasResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&locResp); err != nil {
		return LocationAreasResponse{}, err
	}
	return locResp, nil
}
