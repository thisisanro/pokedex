package main

import (
	"fmt"
	"os"

	"github.com/thisisanro/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
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
	}
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandMap(c *config) error {
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

func commandMapB(c *config) error {
	var url string
	if c.previousURL == nil {
		fmt.Println("you are on the first page")
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
