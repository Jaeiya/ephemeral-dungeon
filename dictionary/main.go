package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jaeiya/monster/dictionary/utils"
)

const filePath = "./dictionary.db"

type FileStorage struct {
}

func (fs FileStorage) Save(b []byte) error {
	fh, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file::%w", err)
	}

	defer fh.Close()

	_, err = fh.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write to file::%w", err)
	}

	return nil
}

func (fs FileStorage) ReadAll() ([]byte, error) {
	fh, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file::%w", err)
	}
	defer fh.Close()
	return io.ReadAll(fh)
}

func displayMenu(r *bufio.Reader, dict Dictionary) {
	for {
		choice, err := utils.PromptMenu(utils.MenuOptions{
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
			utils.PrintError(err)
			utils.PromptBackToMenu(r)
			continue
		}

		switch choice {
		case 1:
			input := strings.ToLower(utils.PromptInput("Add Word", r))

			if !isValidInput(input) {
				utils.PrintError(fmt.Errorf("'%s' is not a valid word string", input))
				break
			}

			if len(input) == 0 {
				utils.PrintError(fmt.Errorf("empty input not allowed"))
				break
			}

			success, err := dict.AddWords(input)
			if err != nil {
				utils.PrintError(err)
			} else if !success {
				utils.PrintError(fmt.Errorf("one or all of '%s' already exists", input))
			}

		case 2:
			input := strings.ToLower(utils.PromptInput("Append Words", r))
			success, err := dict.AppendWord(input)
			if err != nil {
				utils.PrintError(err)
			}
			if !success {
				utils.PrintError(
					fmt.Errorf("one or more of the words to append already exists"),
				)
			}

		case 3:
			input := strings.ToLower(utils.PromptInput("View Word", r))
			wordID, exists := dict.wordMap[input]
			if !exists {
				utils.PrintError(fmt.Errorf("could not find '%s'", input))
				break
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

		case 4:
			for key, val := range dict.wordMap {
				fmt.Println(val, key)
			}

		case 5:
			return
		}

		utils.PromptBackToMenu(r)
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

func main() {
	dict, err := NewDictionary(FileStorage{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n%+v\n%+v\n", hex.EncodeToString(dict.hash), dict.length, dict.wordMap)
	success, err := dict.AppendWord("stone rock")
	if err != nil {
		panic(err)
	}

	if !success {
		panic("tried to add a word that already exists")
	}

	fmt.Println(dict.wordMap)

	// var err error

	// fh, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	// if err != nil {
	// 	panic(err)
	// }
	// defer fh.Close()

	// fileStats, err := fh.Stat()
	// if err != nil {
	// 	panic(err)
	// }

	// var isLoaded bool

	// if store, isLoaded = loadDict(fh, fileStats.Size()); !isLoaded {
	// 	// Write new dictionary file header
	// 	createHeader(fh, hasher, dictLength)
	// 	store = map[string]uint16{}
	// }

	// r := bufio.NewReader(os.Stdin)
	// displayMenu(r, fh)
}
