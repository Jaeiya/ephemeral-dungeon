package lib

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/jaeiya/monster/CLI/shared"
)

func DisplayMenu(r *bufio.Reader, dict *Dictionary) {
	for {
		choice, err := shared.PromptMenu(shared.MenuOptions{
			Title: "Dictionary Config",
			Items: []string{
				"Add Word",
				"Append Word",
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
			if err := promptAddWord(r, dict); err != nil {
				shared.PrintError(err)
			}

		case 2:
			if err := promptAppendWord(r, dict); err != nil {
				shared.PrintError(err)
			}

		case 3:
			if err := printWord(r, dict); err != nil {
				shared.PrintError(err)
			}

		case 4:
			viewWordMap(dict)

		case 5:
			return
		}

		shared.PromptBackToMenu(r)
	}
}

func promptAddWord(r *bufio.Reader, dict *Dictionary) error {
	input := strings.ToLower(shared.PromptInput("Add Word", r))

	if !isValidInput(input) {
		return fmt.Errorf("'%s' is not a valid word string", input)
	}

	if len(input) == 0 {
		return fmt.Errorf("empty input not allowed")
	}

	success, err := dict.AddWords(input)
	if err != nil {
		return err
	} else if !success {
		return fmt.Errorf("one or all of '%s' already exists", input)
	}

	return nil
}

func promptAppendWord(r *bufio.Reader, dict *Dictionary) error {
	input := strings.ToLower(shared.PromptInput("Append Words", r))
	success, err := dict.AppendWord(input)

	if err != nil {
		return err
	} else if !success {
		return fmt.Errorf("one or more of the words to append already exists")
	}

	return nil
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
	for key, val := range dict.wordMap {
		fmt.Println(val, key)
	}
}

// isValidInput returns false for all non-lowercase english letters or spaces
func isValidInput(input string) bool {
	for _, r := range input {
		if r != 32 && (r < 97 || r > 122) {
			return false
		}
	}
	return true
}
