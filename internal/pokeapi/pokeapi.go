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

type PokemonNamesResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
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
	if resp.StatusCode != http.StatusOK {
		return locResp, fmt.Errorf("bad response: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return locResp, err
	}
	cache.Add(url, data)

	// return new response
	if err := json.Unmarshal(data, &locResp); err != nil {
		return locResp, err
	}
	return locResp, nil
}

func GetAreaPokemonNames(url string) ([]string, error) {
	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}

	var pokeNames []string
	var pokeResp PokemonNamesResponse

	// return names if they're in cache
	val, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(val, &pokeResp); err != nil {
			return nil, err
		}
		pokeNames = makeSliceOfNames(pokeResp)
		return pokeNames, nil
	}

	// make new request, add response to cache
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, data)

	//return slice of new names
	if err := json.Unmarshal(data, &pokeResp); err != nil {
		return nil, err
	}
	pokeNames = makeSliceOfNames(pokeResp)
	return pokeNames, nil
}

func makeSliceOfNames(pnr PokemonNamesResponse) []string {
	names := make([]string, 0, len(pnr.PokemonEncounters))
	for _, p := range pnr.PokemonEncounters {
		names = append(names, p.Pokemon.Name)
	}
	return names
}
