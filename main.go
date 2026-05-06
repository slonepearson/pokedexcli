package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/commands"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			err := scanner.Err()
			if err != nil {
				fmt.Printf("Error while scanning Stdin: %v", err)
				break
			}
		}
		userInput := scanner.Text()
		err := commands.LookupCommand(userInput)
		if err != nil {
			fmt.Println(err)
		}
	}
}
