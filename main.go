package main

import "fmt"

// Decorator Pattern

type DB interface {
	StoreToDB(string) error
}

type Store struct{}

func (s *Store) StoreToDB(value string) error {
	fmt.Println("Stored to DB", value)
	return nil
}

func myFunction(db DB) function {
	return func (s string) {
		fmt.Println(s)
		db.StoreToDB(s)
	}
}

func main() {
	s := &Store{}
	Execute(myFunction(s))
}

// third party function
type function func(string)

func Execute(fn function) {
	fn("FOO BAR BAZ")
}
