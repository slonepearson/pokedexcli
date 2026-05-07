package commands

import (
	"fmt"
	"io"
	"os"
	"pokedexcli/internal/pokeapi"
	"strings"
)

// each callback receives an io.Writer(os.Stdout) for easier testing
type cliCommand struct {
	name        string
	description string
	callback    func(io.Writer, string) error
}

// Keeps tracks of the previous and next urls return by api call
// to pagenate through the locations
type config struct {
	previous string
	next     string
}

func (c *config) updateNext(next *string) {
	if next != nil {
		c.next = *next
	} else {
		c.next = ""
	}
}

func (c *config) updatePreviuos(previous *string) {
	if previous != nil {
		c.previous = *previous
	} else {
		c.previous = ""
	}
}

var configPtr *config = &config{}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)
	return words
}

func parseInput(input string) (command string, params string, err error) {
	if strings.TrimSpace(input) == "" {
		return "", "", fmt.Errorf("No command: type 'help' to see the supported commands")
	}
	inputs := cleanInput(input)
	command = inputs[0]
	if len(inputs) > 1 {
		params = inputs[1]
	}
	if len(inputs) > 2 {
		for _, param := range inputs[2:] {
			params += fmt.Sprintf(" %s", param)
		}
	}
	return command, params, nil
}

func LookupCommand(input string, w io.Writer) error {
	command, params, err := parseInput(input)
	if err != nil {
		return err
	}
	supportedCommands := getCommands()
	handler, ok := supportedCommands[command]
	if ok {
		return handler.callback(w, params)
	} else {
		return fmt.Errorf("Unknown command: type 'help' to see the supported commands")
	}
}

// Registry of all supported commands
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas in the Pokemon world",
			callback:    commandMapB,
		},
	}
}

func commandExit(w io.Writer, params string) error {
	fmt.Fprintln(w, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(w io.Writer, params string) error {
	fmt.Fprintln(w, "Welcome to the Pokedex!")
	fmt.Fprint(w, "Usage:\n\n")
	for _, command := range getCommands() {
		fmt.Fprintf(w, "%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(w io.Writer, params string) error {
	url := ""
	if configPtr.next != "" {
		url = configPtr.next
	}

	locationAreas, err := pokeapi.GetLocationAreas(url)

	if err != nil {
		return err
	}

	configPtr.updateNext(locationAreas.Next)
	configPtr.updatePreviuos(locationAreas.Previous)

	for _, result := range locationAreas.Results {
		fmt.Fprintln(w, result.Name)
	}
	return nil
}

func commandMapB(w io.Writer, params string) error {
	url := configPtr.previous
	if url == "" {
		return fmt.Errorf("you're on the first page")
	}

	locationAreas, err := pokeapi.GetLocationAreas(url)

	if err != nil {
		return err
	}

	configPtr.updatePreviuos(locationAreas.Previous)
	configPtr.updateNext(locationAreas.Next)

	for _, result := range locationAreas.Results {
		fmt.Fprintln(w, result.Name)
	}
	return nil
}
