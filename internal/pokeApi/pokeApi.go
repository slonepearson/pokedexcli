package pokeApi

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	res, err := http.Get(url)
	if err != nil {
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
