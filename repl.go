package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/thisisanro/pokedex/internal/pokeapi"
)

type config struct {
	nextURL     *string
	previousURL *string
	pokedex     map[string]pokeapi.PokemonDetails
}

var baseConfig = config{
	pokedex: make(map[string]pokeapi.PokemonDetails),
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			return
		}
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		inputCommand := words[0]

		command, exists := getCommands()[inputCommand]
		if exists {
			err := command.callback(&baseConfig, words[1:])
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

// split user's input by spaces into slice of lowercase strings
func cleanInput(text string) []string {
	lowString := strings.ToLower(text)
	result := strings.Fields(lowString)
	return result
}
