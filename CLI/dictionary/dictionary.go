package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"maps"
	"slices"
	"strings"
)

const (
	// dictLength (BigEndian uint16) + CRC32 + \n
	headerSize = 2 + 4 + 1
	newLineIdx = headerSize - 1
)

type DictStorage interface {
	ReadAll() ([]byte, error)
	Save(data []byte) error
}

type FileInfo struct {
	dictLength uint16
	hash       []byte
}

type Dictionary struct {
	wordMap  map[string]uint16
	length   uint16
	store    DictStorage
	checksum uint32
}

func NewDictionary(store DictStorage) (*Dictionary, error) {
	fileData, err := store.ReadAll()
	if err != nil {
		return &Dictionary{}, fmt.Errorf("error reading file::%w", err)
	}

	if len(fileData) == 0 {
		dict := &Dictionary{}
		newHeader := make([]byte, 0, headerSize)
		newHeader = binary.BigEndian.AppendUint16(newHeader, 0)

		checksum := crc32.ChecksumIEEE([]byte{})
		binary.BigEndian.AppendUint32(newHeader, checksum)
		newHeader = append(newHeader, '\n')

		if err := store.Save(newHeader); err != nil {
			return dict, fmt.Errorf("failed to save new dictionary::%w", err)
		}

		dict.length = 0
		dict.checksum = checksum
		dict.store = store
		dict.wordMap = map[string]uint16{}
		return dict, nil
	}

	header := fileData[:headerSize]
	content := fileData[headerSize:]

	if header[newLineIdx] != '\n' {
		return &Dictionary{}, fmt.Errorf("invalid header::missing termination")
	}

	dict := &Dictionary{
		length:   binary.BigEndian.Uint16(header[0:2]),
		checksum: binary.BigEndian.Uint32(header[2:newLineIdx]),
		wordMap:  map[string]uint16{},
		store:    store,
	}

	if !dict.isValidHash(content) {
		return &Dictionary{}, fmt.Errorf("data integrity check failed")
	}

	for data := range bytes.SplitSeq(content, []byte{'\n'}) {
		if len(data) > 0 {
			id := binary.BigEndian.Uint16(data[:2])
			for word := range strings.SplitSeq(string(data[2:]), " ") {
				dict.wordMap[word] = id
			}
		}
	}

	return dict, nil
}

// AddWords splits the wordInput by space character into a word
// slice and adds them all as a single dictionary entry, where
// they all reference the same ID.
//
// 🟡 Returns false if any words in the wordInput already exist
// in the dictionary.
func (dict *Dictionary) AddWords(wordInput string) (bool, error) {
	words := strings.Fields(wordInput)

	for _, w := range words {
		if _, exists := dict.wordMap[w]; exists {
			return false, nil
		}
	}

	dict.length += 1
	for _, w := range words {
		dict.wordMap[w] = dict.length
	}

	if err := dict.save(); err != nil {
		return false, fmt.Errorf("failed to save words::%w", err)
	}

	return true, nil
}

// AppendWord splits the wordInput by space character into a word
// slice. All words found after the first, are appended to the
// id of the first.
//
// 🟠 If the first word does not exist in the dictionary, an
// error will be thrown.
//
// 🟡 Returns false if any of the words to append, already
// exist in the dictionary.
func (dict *Dictionary) AppendWord(wordInput string) (bool, error) {
	words := strings.Fields(strings.ToLower(wordInput))

	wordID, exists := dict.wordMap[words[0]]
	if !exists {
		return false, fmt.Errorf("first word must exist in dictionary to append to")
	}

	if len(words) < 2 {
		return false, fmt.Errorf("missing words to append to '%s'", words[0])
	}

	wordsToAppend := words[1:]

	for _, w := range wordsToAppend {
		if _, exists := dict.wordMap[w]; exists {
			return false, nil
		}
		dict.wordMap[w] = uint16(wordID)
	}

	if err := dict.save(); err != nil {
		return false, fmt.Errorf("failed to save appended words::%w", err)
	}

	return true, nil
}

// save writes a header and dictionary data through the
// storage interface.
//
// 🔵 The header is made up of a BigEndian uint16
// dictionary length and a SHA-256 hash of the
// dictionary data.
//
// 🔵 The hash is used to check for data corruption
// or tampering.
func (dict *Dictionary) save() error {
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, dict.length)

	if len(dict.wordMap) == 0 {
		buf := [headerSize]byte{}

		bufSlice := buf[:0]
		bufSlice = binary.BigEndian.AppendUint16(bufSlice, 0)
		bufSlice = binary.BigEndian.AppendUint32(bufSlice, crc32.ChecksumIEEE(nil))
		bufSlice = append(bufSlice, '\n')

		return dict.store.Save(bufSlice)
	}

	reverseDict := map[uint16][]string{}
	for k, v := range dict.wordMap {
		reverseDict[v] = append(reverseDict[v], k)
	}

	// Force sequential order of lines in file
	keys := slices.Sorted(maps.Keys(reverseDict))

	buf := bytes.Buffer{}
	idBytes := make([]byte, 2)
	buf.Grow(len(dict.wordMap) * 10)

	for _, k := range keys {
		binary.BigEndian.PutUint16(idBytes, k)
		buf.Write(idBytes)

		v := reverseDict[k]
		slices.Sort(v) // Force alphabetical order of synonyms
		buf.WriteString(strings.Join(v, " "))
		buf.WriteByte('\n')
	}

	bufSize := headerSize + buf.Len()
	fileBuf := make([]byte, 0, bufSize)

	fileBuf = binary.BigEndian.AppendUint16(fileBuf, dict.length)
	fileBuf = binary.BigEndian.AppendUint32(fileBuf, crc32.ChecksumIEEE(buf.Bytes()))
	fileBuf = append(fileBuf, '\n')
	fileBuf = append(fileBuf, buf.Bytes()...)

	return dict.store.Save(fileBuf)
}

func (dict *Dictionary) isValidHash(data []byte) bool {
	return dict.checksum == crc32.ChecksumIEEE(data)
}
