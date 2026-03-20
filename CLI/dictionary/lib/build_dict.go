package lib

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/jaeiya/monster/CLI/shared"
	"github.com/pelletier/go-toml/v2"
)

var TomlDictionary struct {
	Synonyms [][]string `toml:"synonyms"`
}

// dictLength (BigEndian uint16) + CRC32 + \n
const _headerSize = 2 + 4 + 1

var _dbStore = NewFileStorage("./dictionary.db")

func BuildDict(store DictStorage) error {
	data, err := store.ReadAll()
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("specified dictionary file is empty")
	}

	err = toml.Unmarshal(data, &TomlDictionary)
	if err != nil {
		return err
	}

	mapLen := 0
	for i := range TomlDictionary.Synonyms {
		mapLen += len(TomlDictionary.Synonyms[i])
	}

	var length uint16
	wordMap := make(map[string]uint16, mapLen)

	for _, words := range TomlDictionary.Synonyms {
		length++
		for _, word := range words {
			word = strings.ToLower(word)

			if _, exists := wordMap[word]; exists {
				return fmt.Errorf("remove duplicate word: '%s'", word)
			}
			if strings.ContainsRune(word, ' ') {
				return fmt.Errorf("words cannot contain spaces: '%s'", word)
			}
			if !shared.HasLowercaseOnly(word) {
				return fmt.Errorf("word contains invalid characters: '%s'", word)
			}
			wordMap[word] = length
		}
	}

	return save(_dbStore, length, wordMap)
}

// save writes a header and dictionary data through the
// storage interface.
//
// 🔵 The header is made up of a BigEndian uint16
// dictionary length and a CRC32 of the
// dictionary data.
//
// 🔵 The CRC32 is used to check for data corruption
func save(store DictStorage, length uint16, wordMap map[string]uint16) error {
	if len(wordMap) == 0 {
		fileBuf, _ := createHeader(0, nil)
		return store.Save(fileBuf)
	}

	reverseDict := map[uint16][]string{}
	for k, v := range wordMap {
		reverseDict[v] = append(reverseDict[v], k)
	}

	buf := bytes.Buffer{}
	idBytes := make([]byte, 2)
	buf.Grow(len(wordMap) * 10)

	for k, v := range reverseDict {
		binary.BigEndian.PutUint16(idBytes, k)
		buf.Write(idBytes)
		buf.WriteString(strings.Join(v, " "))
		buf.WriteByte('\n')
	}

	data := buf.Bytes()
	header, _ := createHeader(length, data)

	w, err := store.NewWriter()
	if err != nil {
		return err
	}
	defer w.Close()

	bw := bufio.NewWriter(w)
	bw.Write(header)
	bw.Write(data)
	return bw.Flush()
}

func createHeader(length uint16, data []byte) ([]byte, uint32) {
	checksum := crc32.ChecksumIEEE(data)

	header := make([]byte, 0, _headerSize)
	header = binary.BigEndian.AppendUint16(header, length)
	header = binary.BigEndian.AppendUint32(header, checksum)
	header = append(header, '\n')

	return header, checksum
}
