package lib

import (
	"bufio"
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/jaeiya/monster/CLI/shared"
)

const (
	dbPath   = "./dictionary.db"
	tomlPath = "lib/dictionary/dictionary.toml"
)

func DisplayMenu(r *bufio.Reader) error {
	dict, err := NewDictionary(NewFileStorage(dbPath))
	if err != nil {
		panic(err)
	}

	for {
		choice, err := shared.PromptMenu(shared.MenuOptions{
			Title: "Dictionary Config",
			Items: []string{
				"Build",
				"View Word",
				"View All",
			},
			SoftExit: true,
		}, r)
		if err != nil {
			shared.PrintError(err)
			shared.PromptBackToMenu(r)
			continue
		}

		switch choice {

		case 1:
			err := BuildDict(NewFileStorage(tomlPath))
			if err != nil {
				shared.PrintError(err)
				break
			}

			dict, err = NewDictionary(NewFileStorage(dbPath))
			if err != nil {
				shared.PrintError(err)
				break
			}

			fmt.Print(shared.ColorString("\n  ;c;Dictionary Re-Built!"))

		case 2:
			if err := printWord(r, dict); err != nil {
				shared.PrintError(err)
			}

		case 3:
			viewWordMap(dict)

		case 4:
			return nil
		}

		shared.PromptBackToMenu(r)
	}
}

func printWord(r *bufio.Reader, dict *Dictionary) error {
	input := strings.ToLower(shared.PromptInput("View Word", r))
	wordID, exists := dict.wordMap[input]
	if !exists {
		return fmt.Errorf("could not find '%s'", input)
	}

	words := []string{}
	for word, dID := range dict.wordMap {
		if dID == wordID && word != input {
			words = append(words, word)
		}
	}

	fmt.Printf("\n %d %s\n", wordID, input)
	for _, w := range words {
		fmt.Printf(" %d %s\n", wordID, w)
	}

	return nil
}

func viewWordMap(dict *Dictionary) {
	type entry struct {
		word string
		id   uint16
	}

	items := make([]entry, 0, len(dict.wordMap))
	for k, v := range dict.wordMap {
		items = append(items, entry{k, v})
	}

	slices.SortFunc(items, func(a, b entry) int {
		return cmp.Compare(a.id, b.id)
	})

	fmt.Println()
	for _, item := range items {
		fmt.Printf("    %d %s\n", item.id, item.word)
	}
}
