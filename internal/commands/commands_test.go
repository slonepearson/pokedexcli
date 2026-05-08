package commands

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
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

func TestParseInput(t *testing.T) {
	type expected struct {
		command string
		args    []string
		err     error
	}

	cases := []struct {
		input    string
		expected expected
	}{
		{
			input:    "Help me",
			expected: expected{"help", []string{"me"}, nil},
		},
		{
			input:    "exit",
			expected: expected{"exit", []string{}, nil},
		},
		{
			input:    "",
			expected: expected{"", []string{}, errors.New("No command: type 'help' to see the supported commands")},
		},
		{
			input:    "map over the locations",
			expected: expected{"map", []string{"over", "the", "location"}, nil},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test %v", i+1), func(t *testing.T) {
			actualCommand, actualArgs, actualErr := parseInput(c.input)

			if c.expected.err != nil && actualErr == nil {
				t.Errorf("Expected an error got <nil>")
			}

			if c.expected.err == nil && actualErr != nil {
				t.Errorf("Didn't expect an error gor: %v", actualErr)
			}

			if actualCommand != c.expected.command {
				t.Errorf("Expected: %v, Got: %v", c.expected.command, actualCommand)
			}
			if len(actualArgs) != len(c.expected.args) {
				t.Errorf("Expected: %#v, Got: %#v", c.expected.args, actualArgs)
			}
		})
	}
}

func TestHelpCommand(t *testing.T) {
	cases := []struct {
		name      string
		command   string
		toContain []string
	}{
		{
			name:      "commandHelp",
			command:   "help",
			toContain: []string{"help", "exit", "map", "mapb"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := LookupCommand(tc.command, &buf)

			if err != nil {
				t.Fatalf("Expected no errors, got %v", err)
			}

			for _, expected := range tc.toContain {
				if !strings.Contains(buf.String(), expected) {
					t.Fatalf("Expected usage to show command: %s\n Got %s ", expected, buf.String())
				}
			}

		})
	}
}

func TestCommandMap(t *testing.T) {
	cases := []struct {
		name          string
		command       string
		expectedCount int
	}{
		{
			name:          "commandMap",
			command:       "map",
			expectedCount: 20,
		},
		{
			name:          "commandMapNext",
			command:       "map",
			expectedCount: 40,
		},
	}
	var buf bytes.Buffer

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			err := LookupCommand(tc.command, &buf)

			if err != nil {
				t.Fatalf("Expected no errors, got %v", err)
			}

			s := buf.String()
			actual := strings.Split(strings.Trim(s, "\n"), "\n")
			if actualCount := len(actual); tc.expectedCount != actualCount {
				t.Fatalf("Expected count: %d, Actual count: %d", tc.expectedCount, actualCount)
			}
		})
	}
}

func TestCommandMapB(t *testing.T) {
	cases := []struct {
		name          string
		command       string
		expectedCount int
		wantErr       bool
	}{
		{
			name:    "firstPageCommandB",
			command: "mapb",
			wantErr: true,
		},
		{
			name:    "commandMapB",
			command: "mapb",
		},
	}

	configPtr.previous = ""
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			var expectedBuf bytes.Buffer

			if tc.command == "mapb" && !tc.wantErr {
				LookupCommand("map", &expectedBuf)
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

			s := buf.String()
			e := buf.String()
			actual := strings.Split(strings.Trim(s, "\n"), "\n")
			expected := strings.Split(strings.Trim(e, "\n"), "\n")

			if !slices.Equal(actual, expected) {
				t.Fatalf("Expected the result to equal the first page result")
			}

		})
	}
}
