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

const interval = 20 * time.Second

var cache = pokecache.NewCache(interval)

type locationAreas struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type pokemonEncounters struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// GET request to the Location Area Endpoint with cache
func GetLocationAreaEndpoint(url string) ([]byte, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}

	if data, exists := cache.Get(url); exists {
		return data, nil
	}

	client := http.Client{Timeout: 5 * time.Second}

	res, err := client.Get(url)
	if err != nil {
		if os.IsTimeout(err) {
			return nil, fmt.Errorf("The request timed out: %v", err)
		}
		return nil, fmt.Errorf("network error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok GET request: %v", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}

	cache.Add(url, data)
	return data, nil
}

func GetAreas(url string) (locationAreas, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}

	data, err := GetLocationAreaEndpoint(url)
	if err != nil {
		return locationAreas{}, err
	}
	var locations locationAreas
	if err = json.Unmarshal(data, &locations); err != nil {
		return locationAreas{}, fmt.Errorf("issue parsing json: %v", err)
	}
	return locations, err
}

func FindPokemon(url string, areaName string) (pokemonEncounters, error) {
	url = fmt.Sprintf("%s/%v", url, areaName)
	data, err := GetLocationAreaEndpoint(url)
	if err != nil {
		return pokemonEncounters{}, err
	}
	var encounters pokemonEncounters
	if err := json.Unmarshal(data, &encounters); err != nil {
		return pokemonEncounters{}, err
	}
	return encounters, nil
}
