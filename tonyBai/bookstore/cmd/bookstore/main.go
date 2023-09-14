package main

import (
	_ "bookstore/internal/store"
	"bookstore/store/factory"
)

func main() {
	s, err := factory.New("mem")
	if err != nil {
		panic(err)
	}

}
