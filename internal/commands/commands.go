package commands

import (
	"fmt"
	"os"
	"pokedexcli/internal/pokeapi"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

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

func LookupCommand(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("No command: type 'help' to see the supported commands")
	}
	command := cleanInput(input)[0]
	supportedCommands := getCommands()
	handler, ok := supportedCommands[command]
	if ok {
		return handler.callback()
	} else {
		return fmt.Errorf("Unknown command: type 'help' to see the supported commands")
	}
}

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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap() error {
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
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapB() error {
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
		fmt.Println(result.Name)
	}
	return nil
}
