package main

import (
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Début")

	users := []struct {
		ID   int
		Name string
		Age  int
	}{
		{1, "Alice", 25},
		{2, "Bob", 32},
		{3, "Charlie", 41},
	}

	fmt.Println("Contenu :")
	for _, u := range users {
		fmt.Printf("ID=%d Name=%s Age=%d\n", u.ID, u.Name, u.Age)
	}

	fmt.Println("Fin")
}