package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pokedexcli/internal/pokecache"
	"time"
)

var cache = pokecache.NewCache(20 * time.Second)

type areaLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(url string) (areaLocations, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}

	if data, exists := cache.Get(url); exists {
		var locations areaLocations
		err := json.Unmarshal(data, &locations)
		if err != nil {
			return areaLocations{}, fmt.Errorf("issue parsing json: %v", err)
		}
		return locations, nil
	}

	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Get(url)
	if err != nil {
		if os.IsTimeout(err) {
			return areaLocations{}, fmt.Errorf("The request timed out: %v", err)
		}
		return areaLocations{}, fmt.Errorf("network error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return areaLocations{}, fmt.Errorf("non ok GET request: %v", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return areaLocations{}, fmt.Errorf("error reading data: %v", err)
	}

	cache.Add(url, data)

	var locations areaLocations
	if err = json.Unmarshal(data, &locations); err != nil {
		return areaLocations{}, fmt.Errorf("issue parsing json: %v", err)
	}

	return locations, nil
}
