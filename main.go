package main

import (
	"fmt"

	"github.com/NickDubelman/azamat"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// Create todos table
	createTable := `
		CREATE TABLE todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL
		)
	`

	if _, err = db.Exec(createTable); err != nil {
		panic(err)
	}

	type Todo struct {
		ID    int
		Title string
	}

	TodoTable := azamat.Table[Todo]{
		Name:    "todos",
		Columns: []string{"id", "title"},
	}

	// Insert
	todoTitle := "assist Borat"
	insert := TodoTable.Insert().Columns("title").Values(todoTitle)
	result, err := insert.Run(db)
	if err != nil {
		panic(err)
	}
	todoID, _ := result.LastInsertId()
	fmt.Printf("Inserted todo: id=%d, title=%s\n", todoID, todoTitle)

	// Query all
	todos, err := TodoTable.GetAll(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query all:", todos)

	// Query by ID
	todo, err := TodoTable.GetByID(db, int(todoID))
	if err != nil {
		panic(err)
	}

	fmt.Println("Query by ID:", todo)

	todoTitle = "no longer friends with borat"
	update := TodoTable.Update().Set("title", todoTitle).Where("id = ?", todoID)
	if _, err := update.Run(db); err != nil {
		panic(err)
	}

	fmt.Printf("Updated todo: id=%d, title=%s\n", todoID, todoTitle)

	todo, err = TodoTable.GetByID(db, int(todoID))
	if err != nil {
		panic(err)
	}

	fmt.Println("Query by ID:", todos[0])

	delete := TodoTable.Delete().Where("id = ?", todoID)
	if _, err := delete.Run(db); err != nil {
		panic(err)
	}

	fmt.Printf("Deleted todo: id=%d\n", todoID)
}
