package main

import (
	"fmt"
	"go-linq/lib"
)

func main() {
	storage := lib.NewStorage()

	jobs := []*lib.Job{{Name: "student"}}

	person := lib.Person{
		Name:    "Test",
		Surname: "Test",
		Jobs:    jobs,
	}

	p, _ := storage.Store(&person)
	addedPerson := p.(*lib.Person)

	fmt.Println(addedPerson)
	fmt.Println(storage.Load(addedPerson.Id, lib.Person{}))

	out := storage.Where(lib.Person{}, func(elem interface{}) bool {
		p := elem.(*lib.Person)
		return p.Surname == "Test"
	})

	fmt.Println(out)

	fmt.Println(storage.Load(1, lib.Job{}))
}
