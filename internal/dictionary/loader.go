package dictionary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"strings"
)

const (
	// dictLength (BigEndian uint16) + CRC32 + \n
	headerSize = 2 + 4 + 1
	newLineIdx = headerSize - 1
)

var ErrorDictEmpty = fmt.Errorf("specified dictionary file is empty")

type DictStorage interface {
	NewWriter() (io.WriteCloser, error)
	ReadAll() ([]byte, error)
	Save(data []byte) error
}

type Dictionary struct {
	store    DictStorage
	wordMap  map[string]uint16
	checksum uint32
	length   uint16
}

func NewDictionary(store DictStorage) (*Dictionary, error) {
	fileData, err := store.ReadAll()
	if err != nil {
		return &Dictionary{}, fmt.Errorf("error reading file::%w", err)
	}

	if len(fileData) == 0 {
		return nil, ErrorDictEmpty
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

func (dict *Dictionary) isValidHash(data []byte) bool {
	return dict.checksum == crc32.ChecksumIEEE(data)
}
