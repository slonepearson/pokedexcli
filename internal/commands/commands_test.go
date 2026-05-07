package commands

import (
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
