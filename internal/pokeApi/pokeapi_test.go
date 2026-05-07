package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLocationArea(t *testing.T) {

	cases := []struct {
		name         string
		statusCode   int
		body         string
		wantErr      bool
		wantCount    int
		wantLocation string
	}{
		{
			name:         "success",
			statusCode:   200,
			body:         `{"results":[{"name":"canalave-city-area"}]}`,
			wantErr:      false,
			wantCount:    1,
			wantLocation: "canalave-city-area",
		},
		{
			name:       "server error",
			statusCode: 500,
			body:       ``,
			wantErr:    true,
		},
		{
			name:       "bad json",
			statusCode: 200,
			body:       `{not valid json`,
			wantErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.wantErr && tc.name == "server error" {
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tc.body))
				}
			}))
			defer server.Close()

			locations, err := GetLocationAreas(server.URL)

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
