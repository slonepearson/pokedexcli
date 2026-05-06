package commands

import (
	"fmt"
	"os"
	"pokedexcli/internal/pokeApi"
	"pokedexcli/utils"
	"strings"
)

var configPtr *locationConfig = &locationConfig{}

func LookupCommand(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("No command: type 'help' to see the supported commands")
	}
	command := utils.CleanInput(input)[0]
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
	url := "https://pokeapi.co/api/v2/location-area/"

	if configPtr.Next != "" {
		url = configPtr.Next
	}

	locationAreas, err := pokeApi.GetLocationAreas(url)

	if err != nil {
		return err
	}

	if locationAreas.Next != nil {
		configPtr.Next = *locationAreas.Next
	}
	if locationAreas.Previous != nil {
		configPtr.Previous = *locationAreas.Previous
	}

	for _, result := range locationAreas.Results {
		fmt.Println(result.Name)
	}
	return nil
}
