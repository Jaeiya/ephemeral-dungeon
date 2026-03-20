package dictionary

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStorage implements DictStorage for testing purposes.
type mockStorage struct {
	data     []byte
	readErr  error
	saveErr  error
	writeErr error
}

func (m *mockStorage) ReadAll() ([]byte, error) {
	return m.data, m.readErr
}

func (m *mockStorage) Save(data []byte) error {
	m.data = data
	return m.saveErr
}

func (m *mockStorage) NewWriter() (io.WriteCloser, error) {
	if m.writeErr != nil {
		return nil, m.writeErr
	}
	return &mockWriteCloser{m: m}, nil
}

// mockWriteCloser captures the stream of bytes written via NewWriter.
type mockWriteCloser struct {
	m   *mockStorage
	buf bytes.Buffer
}

func (wc *mockWriteCloser) Write(p []byte) (n int, err error) {
	return wc.buf.Write(p)
}

func (wc *mockWriteCloser) Close() error {
	wc.m.data = wc.buf.Bytes()
	return nil
}

// Helper to generate a valid file buffer (header + content)
func generateValidDictFile(length uint16, content []byte) []byte {
	checksum := crc32.ChecksumIEEE(content)
	header := make([]byte, 0, headerSize)
	header = binary.BigEndian.AppendUint16(header, length)
	header = binary.BigEndian.AppendUint32(header, checksum)
	header = append(header, '\n')
	return append(header, content...)
}

func TestNewDictionary(t *testing.T) {
	t.Run("Fails on empty database", func(t *testing.T) {
		store := &mockStorage{}
		_, err := NewDictionary(store)
		assert.ErrorContains(t, err, "dictionary file is empty")
	})

	t.Run("Fails on read error", func(t *testing.T) {
		store := &mockStorage{readErr: errors.New("disk failure")}
		_, err := NewDictionary(store)

		require.Error(t, err)
		assert.ErrorContains(t, err, "disk failure")
	})

	t.Run("Fails on missing or invalid new line terminator", func(t *testing.T) {
		badHeader := []byte{0, 1, 0, 0, 0, 0, 'X'} // 'X' instead of '\n'
		store := &mockStorage{data: badHeader}
		_, err := NewDictionary(store)

		require.Error(t, err)
		assert.ErrorContains(t, err, "invalid header::missing termination")
	})

	t.Run("Fails on checksum mismatch", func(t *testing.T) {
		content := []byte{0, 1, 'f', 'o', 'o', '\n'}
		// Generate valid file, then tamper with content
		data := generateValidDictFile(1, content)
		data[len(data)-2] = 'x' // Change 'o' to 'x'

		store := &mockStorage{data: data}
		_, err := NewDictionary(store)

		require.Error(t, err)
		assert.ErrorContains(t, err, "data integrity check failed")
	})

	t.Run("Successfully loads existing data", func(t *testing.T) {
		line1 := []byte{0, 1, 'w', '1', ' ', 'w', '2', '\n'}
		line2 := []byte{0, 2, 'w', '3', ' ', 'w', '4', '\n'}

		buf := bytes.Buffer{}
		buf.Grow(len(line1) + len(line2))
		buf.Write(line1)
		buf.Write(line2)

		data := generateValidDictFile(1, buf.Bytes())

		store := &mockStorage{data: data}
		dict, err := NewDictionary(store)

		require.NoError(t, err)
		assert.Equal(t, uint16(1), dict.length)
		assert.Equal(t, uint16(1), dict.wordMap["w1"])
		assert.Equal(t, uint16(1), dict.wordMap["w2"])
		assert.Equal(t, uint16(2), dict.wordMap["w3"])
	})
}

// func TestSaveFormatting(t *testing.T) {
// 	store := &mockStorage{}
// 	dict, err := NewDictionary(store)
// 	require.NoError(t, err)

// 	_, err = dict.AddWords("word1 word2")
// 	require.NoError(t, err)
// 	_, err = dict.AppendWord("word2 synm2 synm3")
// 	require.NoError(t, err)
// 	_, err = dict.DeleteWord("synm2")
// 	require.NoError(t, err)

// 	// Load a fresh dictionary from the saved storage to verify binary output
// 	loadedDict, err := NewDictionary(store)
// 	require.NoError(t, err, "failed to read back saved dictionary")

// 	assert.Equal(t, uint16(2), loadedDict.length)

// 	expectedMap := map[string]uint16{
// 		"word1": 1,
// 		"word2": 2,
// 		"synm3": 2,
// 	}

// 	assert.Equal(t, expectedMap, loadedDict.wordMap)
// }
