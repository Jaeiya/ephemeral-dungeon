package main

import (
	"bufio"
	"os"

	"github.com/jaeiya/monster/internal/dictionary"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	err := dictionary.DisplayMenu(r)
	if err != nil {
		panic(err)
	}
}
