package main

import (
	"bufio"
	"os"

	"github.com/jaeiya/monster/CLI/dictionary/lib"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	err := lib.DisplayMenu(r)
	if err != nil {
		panic(err)
	}
}
