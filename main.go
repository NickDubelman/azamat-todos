package main

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
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
	Author   string
	AuthorID int `db:"authorID"`
}

var TodoTable = azamat.Table[struct {
	ID       int
	Title    string
	AuthorID int `db:"authorID"`
}]{
	Name:    "todos",
	Columns: []string{"id", "title", "authorID"},
	RawSchema: `
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		authorID INTEGER NOT NULL,
		
		FOREIGN KEY (authorID) REFERENCES users(id)
	`,
}

var TodoView = azamat.View[Todo]{
	IDFrom: TodoTable,
	Query: func() sq.SelectBuilder {
		join := fmt.Sprintf(
			"%s ON %s.id = %s.authorID",
			UserTable, UserTable, TodoTable,
		)

		return TodoTable.
			Select().
			Columns("name AS author"). // include author name from UserTable
			Join(join)
	},
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

	// Query all todos
	todos, err := TodoView.GetAll(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query all:", todos)

	// Query specific todo by ID
	todo, err := TodoView.GetByID(db, todoID)
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

	// Query updated todo by its ID
	todo, err = TodoView.GetByID(db, todoID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query by ID:", todo)

	// Delete todo by ID
	delete := TodoTable.Delete().Where("id = ?", todoID)
	if _, err := delete.Run(db); err != nil {
		panic(err)
	}
	fmt.Printf("Deleted todo: id=%d\n", todoID)
}
