package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestCleanInput(t *testing.T) {

	cases := map[string]struct {
		input    string
		expected []string
	}{
		"capitalised":                 {input: "hElLO WORLD", expected: []string{"hello", "world"}},
		"leading spaces	":             {input: "     hello world", expected: []string{"hello", "world"}},
		"trailing spaces":             {input: "hellow world     ", expected: []string{"hello", "world"}},
		"leading and trailing spaces": {input: "  hello  world  ", expected: []string{"hello", "world"}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := cleanInput(c.input)
			if len(actual) != len(c.expected) {
				t.Fatalf("expected: %#v, got: %#v", c.expected, actual)
			}
			for i := range actual {
				word := actual[i]
				expectedWord := actual[i]
				if word != expectedWord {
					t.Fatalf("expected: %#v, got: %#v", expectedWord, word)
				}
			}
		})
	}
}

func TestCommands(t *testing.T) {
	cases := []struct {
		name          string
		command       string
		toContain     []string
		expectedCount int
		wantCount     bool
		wantErr       bool
	}{
		{
			name:      "commandHelp",
			command:   "help",
			toContain: []string{"help", "exit", "map", "mapb"},
			wantCount: false,
		},
		{
			name:          "commandMap",
			command:       "map",
			expectedCount: 20,
			wantCount:     true,
		},
		{
			name:          "commandMapB",
			command:       "mapb",
			expectedCount: 20,
			wantCount:     true,
		},
		{
			name:    "firstPageCommandB",
			command: "mapb",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			if tc.command == "mapb" && !tc.wantErr {
				LookupCommand("map", &buf)
				buf.Reset()
			}

			err := LookupCommand(tc.command, &buf)

			if tc.wantErr && err == nil {
				t.Fatalf("Expected an error got <Nil>")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("Expected no errors, got %v", err)
			}

			if !tc.wantCount {
				for _, expected := range tc.toContain {
					if !strings.Contains(buf.String(), expected) {
						t.Fatalf("Expected usage to show command: %s\n Got %s ", expected, buf.String())
					}
				}
			}

			if tc.wantCount {
				s := buf.String()
				actual := strings.Split(strings.Trim(s, "\n"), "\n")
				if actualCount := len(actual); tc.expectedCount != actualCount {
					t.Fatalf("Expected count: %d, Actual count: %d", tc.expectedCount, actualCount)
				}
			}
		})
	}
}
