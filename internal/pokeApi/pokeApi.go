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

type PokemonInfo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Stats          []struct {
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
}

// GET request to PokeApi endpoints with cache
func MakeGetPokeApi(url string) ([]byte, error) {

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
	defer res.Body.Close()

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

func FindAreas(url string) (locationAreas, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}

	data, err := MakeGetPokeApi(url)
	if err != nil {
		return locationAreas{}, err
	}
	var locations locationAreas
	if err = json.Unmarshal(data, &locations); err != nil {
		return locationAreas{}, fmt.Errorf("issue parsing json: %v", err)
	}
	return locations, err
}

func FindPokemonsByArea(url string, areaName string) (pokemonEncounters, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}
	url = fmt.Sprintf("%v/%v", url, areaName)
	data, err := MakeGetPokeApi(url)
	if err != nil {
		return pokemonEncounters{}, err
	}
	var encounters pokemonEncounters
	if err := json.Unmarshal(data, &encounters); err != nil {
		return pokemonEncounters{}, err
	}
	return encounters, nil
}

func GetPokemonInfo(url string, pokemonName string) (PokemonInfo, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/pokemon"
	}
	url = fmt.Sprintf("%v/%v", url, pokemonName)
	data, err := MakeGetPokeApi(url)
	if err != nil {
		return PokemonInfo{}, err
	}
	var pokemon PokemonInfo
	if err := json.Unmarshal(data, &pokemon); err != nil {
		return PokemonInfo{}, err
	}
	return pokemon, nil
}
