package main

import (
	"bufio"
	"os"
)

const filePath = "./dictionary.db"

func main() {
	var err error

	fh, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	fileStats, err := fh.Stat()
	if err != nil {
		panic(err)
	}

	var isLoaded bool

	if dictionary, isLoaded = loadDict(fh, fileStats.Size()); !isLoaded {
		// Write new dictionary file header
		createHeader(fh, dictLength)
		dictionary = map[string]uint16{}
	}

	r := bufio.NewReader(os.Stdin)
	displayMenu(r, fh)
}
