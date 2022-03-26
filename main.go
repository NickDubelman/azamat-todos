package main

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/NickDubelman/azamat"
	_ "github.com/mattn/go-sqlite3"
)

// First, we define the types that will represent our database entities

type User struct {
	ID   int
	Name string
}

type Todo struct {
	ID       int
	Title    string
	Author   string
	AuthorID int `db:"authorID"`
}

// Next, we define the tables that will store our database entities

// UserTable shows a basic example of a Table. We provide the RawSchema of the table
// so that our code knows how to create the table. This table definition also serves
// as documentation...
var UserTable = azamat.Table[User]{
	Name:    "users",
	Columns: []string{"id", "name"},
	RawSchema: `
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	`,
}

// TodoTable shows an example of a Table that doesn't exactly map to our desired
// representation of an entity. A common example is you want to join with another
// table to include another field that is stored in another table.. below, we will
// use a "View" to accomplish this
var TodoTable = azamat.Table[struct {
	// This struct is anonymous because we don't intend for devs to interact with it
	// directly. Instead, they will interact with Todo's via TodoView
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

// TodoView shows an example of a "View".. Views are useful when we want to represent
// an entity that doesn't precisely map to a single table. With a View, we write a
// custom query that maps to the entity
//
// In this case, we want Todo to include the Author's name, but the TodoTable only
// stores authorID, so we have to join with UserTable to get this:
var TodoView = azamat.View[Todo]{
	IDFrom: TodoTable, // when we GetByID, we need to know which table ID comes from
	Query: func() sq.SelectBuilder {
		// We write a custom query to join with the UserTable and get the author name
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
	db, err := azamat.Connect("sqlite3", ":memory:")
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
	fmt.Printf(
		"Inserted todo: id=%d, title=%s, authorID=%d\n",
		todoID, todoTitle, userID,
	)

	// Query all users
	users, err := UserTable.GetAll(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query all users:", users)

	// Query all todos
	todos, err := TodoView.GetAll(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query all todos:", todos)

	// Query specific user by ID
	user, err := UserTable.GetByID(db, userID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query user by ID:", user)

	// Query specific todo by ID
	todo, err := TodoView.GetByID(db, todoID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query todo by ID:", todo)

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

	todos, err = TodoView.GetAll(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("Query all todos:", todos)

	query := UserTable.Select().Where("name = ?", userName)
	users, err = query.All(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("Custom query (all):", users)

	user, err = query.Only(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("Custom query (only):", user)
}
