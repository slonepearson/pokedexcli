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
	callback    func(io.Writer, []string) error
}

// Keeps tracks of 'previous' and 'next' urls to pagenate through locations
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

func parseInput(input string) (command string, args []string, err error) {
	if strings.TrimSpace(input) == "" {
		return command, args, fmt.Errorf("No command: type 'help' to see supported commands")
	}
	inputs := cleanInput(input)
	command = inputs[0]
	if len(inputs) >= 2 {
		args = inputs[1:]
	}

	return command, args, nil
}

func LookupCommand(input string, w io.Writer) error {
	command, args, err := parseInput(input)
	if err != nil {
		return err
	}
	supportedCommands := getCommands()
	handler, ok := supportedCommands[command]
	if ok {
		return handler.callback(w, args)
	} else {
		return fmt.Errorf("Unknown command: type 'help' to see supported commands")
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
		"explore": {
			name:        "explore <area_name>",
			description: "Display all pokemon listed at a location name",
			callback:    commandExplore,
		},
	}
}

func commandExit(w io.Writer, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many arguments, Usage: exit")
	}
	fmt.Fprintln(w, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(w io.Writer, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many arguments, Usage: help")
	}
	fmt.Fprintln(w, "Welcome to the Pokedex!")
	fmt.Fprint(w, "Usage:\n\n")
	for _, command := range getCommands() {
		fmt.Fprintf(w, "%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(w io.Writer, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many arguments, Usage: map")
	}
	url := ""
	if configPtr.next != "" {
		url = configPtr.next
	}

	locationAreas, err := pokeapi.GetAreas(url)

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

func commandMapB(w io.Writer, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many arguments, Usage: mapb")
	}
	url := configPtr.previous
	if url == "" {
		return fmt.Errorf("you're on the first page")
	}

	locationAreas, err := pokeapi.GetAreas(url)

	if err != nil {
		return err
	}

	configPtr.updatePreviuos(locationAreas.Previous)
	configPtr.updateNext(locationAreas.Next)

	for _, result := range locationAreas.Results {
		fmt.Fprintf(w, "%s\n", result.Name)
	}
	return nil
}

func commandExplore(w io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: explore <location-area-name>")
	}

	encounters, err := pokeapi.FindPokemon("https://pokeapi.co/api/v2/location-area", args[0])
	if err != nil {
		return err
	}

	for _, encounter := range encounters.PokemonEncounters {
		fmt.Fprintf(w, "%s\n", encounter.Pokemon.Name)
	}
	return nil
}
