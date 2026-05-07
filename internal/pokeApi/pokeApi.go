package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

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
		url = "https://pokeapi.co/api/v2/location-area/"
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

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	var locationAreas areaLocations
	if err = decoder.Decode(&locationAreas); err != nil {
		return areaLocations{}, fmt.Errorf("issue parsing json: %v", err)
	}

	return locationAreas, nil
}
