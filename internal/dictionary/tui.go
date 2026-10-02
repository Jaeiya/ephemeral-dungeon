package dictionary

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/jaeiya/monster/internal/utils"
)

const (
	dbPath   = "./data/dictionary.db"
	tomlPath = "./assets/dictionary.toml"
)

type ViewState uint8

const (
	ViewBuild = ViewState(iota)
	ViewWord
	ViewAll
	ViewExit
)

var (
	dbStore   = NewFileStorage(dbPath)
	tomlStore = NewFileStorage(tomlPath)
)

func DisplayMenu(r *bufio.Reader) error {
	dict, err := NewDictionary(NewFileStorage(dbPath))
	if err != nil && !errors.Is(err, ErrorDictEmpty) {
		panic(err)
	}

	for {
		choice, err := utils.PromptMenu(utils.MenuOptions{
			Title: "Dictionary Config",
			Items: []string{
				"Build",
				"View Word",
				"View All",
			},
			SoftExit: true,
		}, r)
		if err != nil {
			utils.PrintError(err)
			utils.PromptBackToMenu(r)
			continue
		}

		switch ViewState(choice) {

		case ViewBuild:
			err := BuildDict(tomlStore, dbStore)
			if err != nil {
				utils.PrintError(err)
				break
			}

			dict, err = NewDictionary(NewFileStorage(dbPath))
			if err != nil {
				utils.PrintError(err)
				break
			}

			fmt.Print(utils.ColorString("\n  ;c;Dictionary Re-Built!"))

		case ViewWord:
			if dict == nil {
				utils.PrintError(ErrorDictEmpty)
				break
			}

			if err := printWord(r, dict); err != nil {
				utils.PrintError(err)
			}

		case ViewAll:
			if dict == nil {
				utils.PrintError(ErrorDictEmpty)
				break
			}

			viewWordMap(dict)

		case ViewExit:
			return nil
		}

		utils.PromptBackToMenu(r)
	}
}

func printWord(r *bufio.Reader, dict *Dictionary) error {
	input := strings.ToLower(utils.PromptInput("View Word", r))
	wordID, exists := dict.WordMap[input]
	if !exists {
		return fmt.Errorf("could not find '%s'", input)
	}

	words := []string{}
	for word, dID := range dict.WordMap {
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

	items := make([]entry, 0, len(dict.WordMap))
	for k, v := range dict.WordMap {
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
