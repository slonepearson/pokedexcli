package commands

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"pokedexcli/internal/pokeapi"
	"strings"
	"time"
)

// Used for storing caught pokemon
var pokedex = map[string]pokeapi.PokemonInfo{}

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
			name:        "explore <area-name>",
			description: "Display all pokemon listed at a location name",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon-name>",
			description: "Catch a desired pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon-name>",
			description: "Inspect a caught pokemon's stats",
			callback:    commandInspect,
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

	locationAreas, err := pokeapi.FindAreas(url)

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

	locationAreas, err := pokeapi.FindAreas(url)

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

	encounters, err := pokeapi.FindPokemonsByArea("", args[0])
	if err != nil {
		return err
	}

	for _, encounter := range encounters.PokemonEncounters {
		fmt.Fprintf(w, "%s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(w io.Writer, args []string) error {

	if len(args) != 1 {
		return fmt.Errorf("usage: catch <pokemon-name>")
	}

	_, ok := pokedex[args[0]]
	if ok {
		return fmt.Errorf("pokemon %s already caught", args[0])
	}

	const k int = 5000
	const minRoll int = 1
	const maxRoll int = 255

	fmt.Fprintf(w, "Throwing a Pokeball at %v...\n", args[0])
	pokemon, err := pokeapi.GetPokemonInfo("", args[0])
	if err != nil {
		return err
	}

	baseExp := pokemon.BaseExperience
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	roll := math.Max(float64(minRoll), float64(r.Intn(maxRoll)))
	threshhold := math.Min(float64(maxRoll), float64(k/baseExp))

	if roll > threshhold {
		fmt.Fprintf(w, "%s escaped!\n", args[0])
		return nil
	}
	pokedex[args[0]] = pokemon
	fmt.Fprintf(w, "%s was caught!\n", args[0])
	return nil
}

func commandInspect(w io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: inspect <caught-pokemon-name>")
	}
	pokemon, ok := pokedex[args[0]]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Fprintf(w, "Name: %v\n", pokemon.Name)
	fmt.Fprintf(w, "Height: %v\n", pokemon.Height)
	fmt.Fprintf(w, "Weight: %v\n", pokemon.Weight)
	fmt.Fprint(w, "Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Fprintf(w, " -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Fprint(w, "Types:\n")
	for _, t := range pokemon.Types {
		fmt.Fprintf(w, " - %v\n", t.Type.Name)
	}
	return nil
}
