package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thisisanro/pokedex/internal/pokecache"
)

const BaseURL = "https://pokeapi.co/api/v2/location-area"

var client = &http.Client{
	Timeout: time.Second * 23,
}

var cache = pokecache.NewCache(time.Second * 5)

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

	var locResp LocationAreasResponse

	// return response if it's in cache
	val, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(val, &locResp); err != nil {
			return LocationAreasResponse{}, err
		}
		return locResp, nil
	}

	// make new request and add to cache
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return locResp, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return locResp, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return locResp, err
	}
	cache.Add(url, data)

	if err := json.Unmarshal(data, &locResp); err != nil {
		return locResp, err
	}
	return locResp, nil
}
