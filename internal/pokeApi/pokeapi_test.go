package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAreas(t *testing.T) {

	cases := []struct {
		name         string
		body         string
		wantErr      bool
		wantCount    int
		wantLocation string
	}{
		{
			name:         "success",
			body:         `{"results":[{"name":"canalave-city-area"}]}`,
			wantErr:      false,
			wantCount:    1,
			wantLocation: "canalave-city-area",
		},
		{
			name:    "network error",
			body:    `{"results":[{"name":"canalave-city-area"}]}`,
			wantErr: true,
		},
		{
			name:    "bad json",
			body:    `{not valid json`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.wantErr && tc.name == "network error" {
					w.WriteHeader(http.StatusBadGateway)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tc.body))
				}
			}))
			defer server.Close()

			locations, err := FindAreas(server.URL)

			if tc.wantErr && err == nil {
				t.Fatalf("expected an error, got nil")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !tc.wantErr && len(locations.Results) != tc.wantCount {
				t.Fatalf("expected %d result, got %d", tc.wantCount, len(locations.Results))
			}
			if !tc.wantErr && locations.Results[0].Name != tc.wantLocation {
				t.Errorf("expected '%s', got %q", tc.wantLocation, locations.Results[0].Name)
			}
		})
	}
}

func TestFindPokemon(t *testing.T) {
	cases := []struct {
		name      string
		path      string
		location  string
		body      string
		wantErr   bool
		wantCount int
		wantName  string
	}{
		{
			name:      "success",
			path:      "/api/v2/location-area",
			location:  "canalave-city-area",
			body:      `{"pokemon_encounters":[{"pokemon":{"name": "staryu"}}]}`,
			wantErr:   false,
			wantCount: 1,
			wantName:  "staryu",
		},
		{
			name:     "network error",
			path:     "/api/v2/location-area",
			location: "canalave-city-area",
			body:     `{"pokemon_encounters":[{"pokemon":{"name": "staryu"}}]}`,
			wantErr:  true,
		},
		{
			name:     "bad json",
			path:     "/api/v2/location-area",
			location: "canalave-city-area",
			body:     `{not valid json`,
			wantErr:  true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.wantErr && tc.name == "network error" {
					w.WriteHeader(http.StatusBadGateway)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tc.body))
					if r.URL.Path != fmt.Sprintf("%v/%v", tc.path, tc.location) {
						t.Fatalf("expected url path: %v", tc.path)
					}
				}
			}))
			defer server.Close()

			encounters, err := FindPokemonsByArea(server.URL+tc.path, tc.location)

			if tc.wantErr && err == nil {
				t.Fatalf("expected an error got <nil>")
			}
			if tc.wantErr {
				return
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected no errors got: %v", err)
			}
			if tc.wantCount != len(encounters.PokemonEncounters) {
				t.Fatalf("expected %d result, got %d", tc.wantCount, len(encounters.PokemonEncounters))
			}
			if tc.wantName != encounters.PokemonEncounters[0].Pokemon.Name {
				t.Fatalf("expected %v, got %v", tc.wantName, encounters.PokemonEncounters[0].Pokemon.Name)
			}
		})
	}
}
