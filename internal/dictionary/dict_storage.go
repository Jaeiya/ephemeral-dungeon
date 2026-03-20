package dictionary

import (
	"fmt"
	"io"
	"os"
)

type FileStorage struct {
	storagePath string
}

func NewFileStorage(path string) FileStorage {
	return FileStorage{path}
}

func (fs FileStorage) Save(b []byte) error {
	fh, err := os.OpenFile(fs.storagePath, os.O_RDWR|os.O_CREATE, 0o644)
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
	fh, err := os.OpenFile(fs.storagePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file::%w", err)
	}
	defer fh.Close()
	return io.ReadAll(fh)
}

func (fs FileStorage) NewWriter() (io.WriteCloser, error) {
	return os.OpenFile(fs.storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
}
