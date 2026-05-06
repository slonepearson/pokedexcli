package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		if strings.TrimSpace(userInput) == "" {
			continue
		}
		cleanedWords := cleanInput(userInput)
		command := cleanedWords[0]

		fmt.Printf("Your command was: %s\n", command)
	}
}
