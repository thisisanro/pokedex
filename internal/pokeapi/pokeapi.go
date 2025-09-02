package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thisisanro/pokedex/internal/pokecache"
)

const BaseURL = "https://pokeapi.co/api/v2"

var client = &http.Client{
	Timeout: time.Second * 23,
}

var cache = pokecache.NewCache(time.Minute * 5)

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

type PokemonDetails struct {
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Name           string `json:"name"`
	ID             int    `json:"id"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
}

// get locations area struct from url
func GetLocations(url string) (LocationAreasResponse, error) {
	var locResp LocationAreasResponse

	data, err := fetchData(url)
	if err != nil {
		return locResp, err
	}

	if err := json.Unmarshal(data, &locResp); err != nil {
		return locResp, err
	}
	return locResp, nil
}

// get pokemon names from location area
func GetAreaPokemonNames(url string) ([]string, error) {
	var pokeNames []string
	var pokeResp PokemonNamesResponse

	data, err := fetchData(url)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &pokeResp); err != nil {
		return nil, err
	}
	pokeNames = makeSliceOfNames(pokeResp)
	return pokeNames, nil
}

func GetPokemon(url string) (PokemonDetails, error) {
	var pokeDetails PokemonDetails

	data, err := fetchData(url)
	if err != nil {
		return pokeDetails, err
	}

	if err := json.Unmarshal(data, &pokeDetails); err != nil {
		return pokeDetails, err
	}
	return pokeDetails, nil
}

func makeSliceOfNames(pnr PokemonNamesResponse) []string {
	names := make([]string, 0, len(pnr.PokemonEncounters))
	for _, p := range pnr.PokemonEncounters {
		names = append(names, p.Pokemon.Name)
	}
	return names
}

func fetchData(url string) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}

	// check cache and return if hit
	if val, ok := cache.Get(url); ok {
		return val, nil
	}

	// make new request, return response data and add it to cache
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
	return data, nil
}
