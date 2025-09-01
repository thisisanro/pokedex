package main

import (
	"fmt"
	"os"

	"github.com/thisisanro/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, args []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Dispalays previous 20 location areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Displays a list of Pokémon found at a location: explore <location-name>",
			callback:    commandExplore,
		},
	}
}

func commandExit(c *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandMap(c *config, args []string) error {
	var url string
	if c.nextURL == nil {
		url = pokeapi.BaseURL
	} else {
		url = *c.nextURL
	}

	locations, err := pokeapi.GetLocations(url)
	if err != nil {
		return err
	}

	c.nextURL, c.previousURL = locations.Next, locations.Previous
	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapB(c *config, args []string) error {
	var url string
	if c.previousURL == nil {
		fmt.Println("You are on the first page")
		return nil
	}

	url = *c.previousURL
	locations, err := pokeapi.GetLocations(url)
	if err != nil {
		return err
	}

	c.nextURL, c.previousURL = locations.Next, locations.Previous

	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandExplore(c *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("explore usage: explore <area_name>")
	}

	area := args[0]
	url := pokeapi.BaseURL + "/" + area
	fmt.Printf("Exploring %v\n", area)
	fmt.Println("Found Pokemon:")

	names, err := pokeapi.GetAreaPokemonNames(url)
	if err != nil {
		return err
	}

	for _, n := range names {
		fmt.Printf("- %s\n", n)
	}
	return nil
}
