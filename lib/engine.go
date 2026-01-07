package lib

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type GameState struct {
	player Character
	cmdMap map[string]func([]string)
}

func NewGame() error {
	scanner := bufio.NewScanner(os.Stdin)
	gs := GameState{}

	gs.cmdMap = map[string]func([]string){
		"look": gs.Look,
		"l":    gs.Look,

		"examine": gs.Examine,
		"x":       gs.Examine,

		"inv":       gs.Inventory,
		"inventory": gs.Inventory,
		"i":         gs.Inventory,
	}

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		text := scanner.Text()

		if !gs.ParseCommand(text) {
			fmt.Printf("You entered: %s\n\n", text)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (g *GameState) ParseCommand(input string) bool {
	inputParts := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")

	if len(inputParts) > 0 {
		if cmd, exists := g.cmdMap[inputParts[0]]; exists {
			cmd(inputParts[1:])
			return true
		}
	}

	return false
}

func (g *GameState) Look(args []string) {
	argLen := len(args)

	if argLen == 0 {
		// Look within general area
		fmt.Printf("Looking around the room I see nothing\n\n")
		return
	}

	switch args[0] {

	// We're looking inside of an inventory
	case "in", "inside", "into":
		object := strings.Join(args[1:], " ")
		// Run a check to see if the object is an inventory
		fmt.Printf("Looking inside %s\n\n", object)
		return

	// We're looking at something specific
	default:
		startIndex := 0
		if args[0] == "at" {
			startIndex = 1
		}

		object := strings.Join(args[startIndex:], " ")
		fmt.Printf("We're examining %s\n\n", object)
		return
	}
}

func (g *GameState) Examine(args []string) {
	if len(args) == 0 {
		fmt.Printf("Nothing to examine\n\n")
		return
	}
	g.Look(slices.Concat([]string{"at"}, args))
}

func (g GameState) Inventory(args []string) {
	fmt.Println("Display the inventory")
}
