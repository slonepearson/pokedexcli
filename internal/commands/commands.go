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
	callback    func(io.Writer) error
}

// Keeps tracks of the previous and next urls return by api call
// to pagenate through the locations
type locationConfig struct {
	Previous string
	Next     string
}

var configPtr *locationConfig = &locationConfig{}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)
	return words
}

func LookupCommand(input string, w io.Writer) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("No command: type 'help' to see the supported commands")
	}
	command := cleanInput(input)[0]
	supportedCommands := getCommands()
	handler, ok := supportedCommands[command]
	if ok {
		return handler.callback(w)
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

func commandExit(w io.Writer) error {
	fmt.Fprintln(w, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(w io.Writer) error {
	fmt.Fprintln(w, "Welcome to the Pokedex!")
	fmt.Fprint(w, "Usage:\n\n")
	for _, command := range getCommands() {
		fmt.Fprintf(w, "%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(w io.Writer) error {
	url := ""
	if configPtr.Next != "" {
		url = configPtr.Next
	}

	locationAreas, err := pokeapi.GetLocationAreas(url)

	if err != nil {
		return err
	}

	if locationAreas.Next != nil {
		configPtr.Next = *locationAreas.Next
	} else {
		configPtr.Next = ""
	}

	if locationAreas.Previous != nil {
		configPtr.Previous = *locationAreas.Previous
	}

	for _, result := range locationAreas.Results {
		fmt.Fprintln(w, result.Name)
	}
	return nil
}

func commandMapB(w io.Writer) error {
	url := configPtr.Previous
	if url == "" {
		return fmt.Errorf("you're on the first page")
	}

	locationAreas, err := pokeapi.GetLocationAreas(url)

	if err != nil {
		return err
	}

	if locationAreas.Previous != nil {
		configPtr.Previous = *locationAreas.Previous
	} else {
		configPtr.Previous = ""
	}

	if locationAreas.Next != nil {
		configPtr.Next = *locationAreas.Next
	}

	for _, result := range locationAreas.Results {
		fmt.Fprintln(w, result.Name)
	}
	return nil
}
