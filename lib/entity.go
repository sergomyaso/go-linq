package lib

type Person struct {
	Id      int
	Name    string
	Surname string
	Jobs    []*Job `linq:"many_to_many"`
}

type Job struct {
	Id     int
	Name   string
	People []*Person `linq:"many_to_many""`
}

type ProjectedPerson struct {
	Id   int
	Name string
}

type GroupByName struct {
	Name string
}
