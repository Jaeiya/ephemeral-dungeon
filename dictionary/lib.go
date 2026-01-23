package main

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/jaeiya/monster/dictionary/utils"
)

const (
	hashPass = "1337420"
	// dictLength (BigEndian uint16) + Sha256 + \n
	headerSize = 2 + 32 + 1
	newLineIdx = 34
)

var (
	hasher     hash.Hash = hmac.New(sha256.New, []byte(hashPass))
	dictLength uint16    = 0
)

type FileInfo struct {
	dictLength uint16
	hash       []byte
}

var dictionary map[string]uint16

func loadDict(f *os.File, fileSize int64) (map[string]uint16, bool) {
	if fileSize == 0 {
		return nil, false
	}

	fileReader := bufio.NewReader(f)

	fileInfo, err := readHeader(fileSize, fileReader)
	if err != nil {
		panic(fmt.Errorf("failed to read header::%w", err))
	}

	tee := io.TeeReader(fileReader, hasher)

	dictLength = fileInfo.dictLength
	dict := make(map[string]uint16, dictLength)
	r := bufio.NewReader(tee)

	for range dictLength {
		idBytes := [2]byte{}
		io.ReadFull(r, idBytes[:])
		id := binary.BigEndian.Uint16(idBytes[:])
		wordPack, err := r.ReadString('\n')
		if err != nil {
			panic(fmt.Errorf("failed to load dict::%w", err))
		}
		wordPack = strings.TrimSpace(wordPack)
		if strings.Contains(wordPack, " ") {
			for w := range strings.FieldsSeq(wordPack) {
				dict[w] = id
			}
		} else {
			dict[wordPack] = id
		}
	}

	if !bytes.Equal(hasher.Sum(nil), fileInfo.hash) {
		panic(fmt.Errorf("dictionary integrity compromised"))
	}

	return dict, true
}

func readHeader(fileSize int64, r *bufio.Reader) (FileInfo, error) {
	fi := FileInfo{}

	if fileSize < headerSize {
		return fi, fmt.Errorf("corrupted or invalid dictionary header")
	}

	header := make([]byte, headerSize)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return fi, fmt.Errorf("failed to read header::%w", err)
	}

	if header[newLineIdx] != '\n' {
		return fi, fmt.Errorf("missing line-termination char")
	}

	fi.dictLength = binary.BigEndian.Uint16(header[0:2])
	fi.hash = slices.Clone(header[2:newLineIdx])

	return fi, nil
}

func displayMenu(r *bufio.Reader, f *os.File) {
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

			if !addWords(f, input) {
				utils.PrintError(fmt.Errorf("one or all of '%s' already exists", input))
			}

		case 2:
			input := strings.ToLower(utils.PromptInput("Append Words", r))
			err := appendWord(input)
			if err != nil {
				utils.PrintError(err)
			}

		case 3:
			input := strings.ToLower(utils.PromptInput("View Word", r))
			wordID, exists := dictionary[input]
			if !exists {
				utils.PrintError(fmt.Errorf("could not find '%s'", input))
				break
			}

			words := []string{}
			for word, dID := range dictionary {
				if dID == wordID && word != input {
					words = append(words, word)
				}
			}

			fmt.Printf("\n %d %s\n", wordID, input)
			for _, w := range words {
				fmt.Printf(" %d %s\n", wordID, w)
			}

		case 4:
			for key, val := range dictionary {
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

// addWords adds all words as synonyms of each other
func addWords(f *os.File, wordInput string) bool {
	words := strings.Fields(wordInput)

	for _, w := range words {
		if _, exists := dictionary[w]; exists {
			return false
		}
	}

	dictLength += 1
	for _, w := range words {
		dictionary[w] = dictLength
	}

	dataSize := 2 + len(wordInput) + 1
	data := make([]byte, dataSize)
	binary.BigEndian.PutUint16(data[:2], dictLength)
	copy(data[2:], wordInput)
	data[dataSize-1] = '\n'
	save(f, data)
	return true
}

func save(f *os.File, data []byte) {
	w := io.MultiWriter(f, hasher)
	if _, err := w.Write(data); err != nil {
		panic(fmt.Errorf("failed to save words::%w", err))
	}
	if err := createHeader(f, dictLength); err != nil {
		panic(fmt.Errorf("failed to save words::%w", err))
	}
}

func createHeader(fh *os.File, length uint16) error {
	fh.Seek(0, io.SeekStart)
	var buf [headerSize]byte
	binary.BigEndian.PutUint16(buf[:2], length)
	copy(buf[2:], hasher.Sum(nil))
	buf[newLineIdx] = '\n'
	if _, err := fh.Write(buf[:]); err != nil {
		return fmt.Errorf("failed to create header::%w", err)
	}
	fh.Seek(0, io.SeekEnd)
	return nil
}

// appendWord appends the specified 'words' to the existing
// specified 'word'
func appendWord(wordInput string) error {
	words := strings.Fields(strings.ToLower(wordInput))

	wordID, exists := dictionary[words[0]]
	if !exists {
		return fmt.Errorf("first word must exist in dictionary to append to")
	}

	if len(words) < 2 {
		return fmt.Errorf("missing words to append to '%s'", words[0])
	}

	for _, w := range words {
		if _, exists := dictionary[w]; exists {
			continue
		}
		dictionary[w] = uint16(wordID)
	}

	return nil
}

func deleteWord() {
}
