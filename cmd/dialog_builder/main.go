package main

import (
	"database/sql"

	"github.com/jaeiya/monster/internal/dialog"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./data/game.db?_pragma=foreign_keys(1)")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	// Prevent "database is locked" errors during heavy writes
	db.SetMaxOpenConns(1)

	err = dialog.BuildDialog(db)
	if err != nil {
		panic(err)
	}
	// questions, err := dialog.BuildDialog()
	// if err != nil {
	// 	panic(err)
	// }
	// _ = questions
	// fmt.Println(questions[0].Answer)
}
