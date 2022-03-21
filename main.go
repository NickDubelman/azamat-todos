package main

import (
	"fmt"

	"github.com/NickDubelman/azamat"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int
	Name string
}

var UserTable = azamat.Table[User]{
	Name:    "users",
	Columns: []string{"id", "name"},
	RawSchema: `
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	`,
}

type Todo struct {
	ID       int
	Title    string
	AuthorID int `db:"authorID"`
}

var TodoTable = azamat.Table[Todo]{
	Name:    "todos",
	Columns: []string{"id", "title", "authorID"},
	RawSchema: `
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		authorID INTEGER NOT NULL,
		
		FOREIGN KEY (authorID) REFERENCES users(id)
	`,
}

func main() {
	// Connect to db
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// Create users table
	if err := UserTable.CreateIfNotExists(db); err != nil {
		panic(err)
	}

	// Create todos table
	if err := TodoTable.CreateIfNotExists(db); err != nil {
		panic(err)
	}

	// Insert a user
	userName := "Azamat"
	insert := UserTable.Insert().Columns("name").Values(userName)
	userID, err := insert.Run(db)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted user: id=%d, name=%s\n", userID, userName)

	// Insert a todo
	todoTitle := "assist Borat"
	insert = TodoTable.
		Insert().
		Columns("title", "authorID").
		Values(todoTitle, userID)

	todoID, err := insert.Run(db)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted todo: id=%d, title=%s\n", todoID, todoTitle)

	// Query all entries
	todos, err := TodoTable.GetAll(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query all:", todos)

	// Query specific entry by ID
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

	// Query updated row by its ID
	todo, err = TodoTable.GetByID(db, todoID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query by ID:", todo)

	// Delete entry by ID
	delete := TodoTable.Delete().Where("id = ?", todoID)
	if _, err := delete.Run(db); err != nil {
		panic(err)
	}
	fmt.Printf("Deleted todo: id=%d\n", todoID)
}
