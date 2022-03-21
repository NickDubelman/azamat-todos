package main

import (
	"fmt"

	"github.com/NickDubelman/azamat"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID    int
	Title string
}

var TodoTable = azamat.Table[Todo]{
	Name:    "todos",
	Columns: []string{"id", "title"},
	RawSchema: `
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL
	`,
}

func main() {
	// Connect to db
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// Create todos table
	if err := TodoTable.CreateIfNotExists(db); err != nil {
		panic(err)
	}

	// Insert
	todoTitle := "assist Borat"
	insert := TodoTable.Insert().Columns("title").Values(todoTitle)
	todoID, err := insert.Run(db)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted todo: id=%d, title=%s\n", todoID, todoTitle)

	// Query all
	todos, err := TodoTable.GetAll(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query all:", todos)

	// Query by ID
	todo, err := TodoTable.GetByID(db, todoID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query by ID:", todo)

	todoTitle = "no longer friends with borat"
	update := TodoTable.
		Update().
		Set("title", todoTitle).
		Where("id = ?", todoID)

	if _, err := update.Run(db); err != nil {
		panic(err)
	}
	fmt.Printf("Updated todo: id=%d, title=%s\n", todoID, todoTitle)

	todo, err = TodoTable.GetByID(db, todoID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query by ID:", todo)

	delete := TodoTable.Delete().Where("id = ?", todoID)
	if _, err := delete.Run(db); err != nil {
		panic(err)
	}
	fmt.Printf("Deleted todo: id=%d\n", todoID)
}
