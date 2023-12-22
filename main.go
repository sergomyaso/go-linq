package main

import (
	"fmt"
	"go-linq/lib"
)

func main() {
	storage := lib.NewStorage()

	job := &lib.Job{
		Name: "student",
	}

	jobs := []*lib.Job{job}

	person := lib.Person{
		Name:    "Test",
		Surname: "Test",
		Jobs:    jobs,
	}

	personSec := lib.Person{
		Name:    "Second",
		Surname: "Test",
		Jobs:    jobs,
	}

	job.People = []*lib.Person{&person, &personSec}

	p, _ := storage.Store(&person)
	addedPerson := p.(*lib.Person)

	fmt.Println(addedPerson)
	fmt.Println(storage.Load(addedPerson.Id, lib.Person{}))

	pipe := lib.NewQueryCmd()

	out := pipe.Where(storage, lib.Person{}, func(elem interface{}) bool {
		p := elem.(*lib.Person)
		return p.Surname == "Test"
	}).Project(storage, lib.Person{}, lib.ProjectedPerson{}).GroupBy(
		storage, lib.Person{}, lib.GroupByName{}, "", func(acc any, elem interface{}) any {
			p := elem.(lib.ProjectedPerson)
			cur := acc.(string)
			return cur + p.Name
		},
	).Result()

	fmt.Println("GROUP_BY", out)

	fmt.Println(storage.Load(1, lib.Person{}))
}
