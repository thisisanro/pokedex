package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"

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
		"catch": {
			name:        "catch",
			description: "Try to catch pokemon by name: catch <name>",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays the name of a Pokemon and prints the name, height, weight, stats and type(s): usage: inspect <name>",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays all the names of the Pokemon the user has caught",
			callback:    commandPokedex,
		},
	}
}

func commandExit(c *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args []string) error {
	fmt.Printf("Pokedex Usage:\n\n")
	commands := make([]string, 0, len(getCommands()))
	for _, command := range getCommands() {
		commands = append(commands, fmt.Sprintf("%s: %s", command.name, command.description))
	}
	sort.Strings(commands)
	for _, c := range commands {
		fmt.Println(c)
	}
	fmt.Println()
	return nil
}

func commandMap(c *config, args []string) error {
	var url string
	if c.nextURL == nil {
		url = pokeapi.BaseURL + "/" + "location-area"
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
	url := pokeapi.BaseURL + "/" + "location-area" + "/" + area
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

func commandCatch(c *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("catch usage: catch <pokemon-name>")
	}

	name := args[0]
	url := pokeapi.BaseURL + "/" + "pokemon" + "/" + name
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := pokeapi.GetPokemon(url)
	if err != nil {
		return err
	}

	chance := catchChance(pokemon.BaseExperience)
	if rand.Intn(100) >= chance {
		fmt.Printf("%s escaped!\n", name)
	} else {
		c.pokedex[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", name)
		fmt.Println("You may now inspect it with the inspect command.")
	}
	return nil
}

func catchChance(exp int) int {
	easy := 80
	normal := 50
	hard := 20

	switch {
	case exp <= 40:
		return easy
	case 40 < exp && exp <= 135:
		return normal
	default:
		return hard
	}
}

func commandInspect(c *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("inspect usage: inspect <pokemon-name>")
	}

	name := args[0]

	p, exists := c.pokedex[name]
	if !exists {
		return fmt.Errorf("you have not caught that pokemon yet")
	}

	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", p.Name, p.Height, p.Weight)
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(c *config, args []string) error {
	fmt.Println("Your Pokedex:")
	if len(c.pokedex) == 0 {
		fmt.Println(" - (pokedex is empty)")
		return nil
	}

	names := make([]string, 0, len(c.pokedex))
	for name := range c.pokedex {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
