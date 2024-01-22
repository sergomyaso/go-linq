package lib

// for many to many testing
type Person struct {
	Id      int
	Name    string
	Surname string
	Jobs    []*Job `linq:"many_to_many"`
}

type Job struct {
	Id     int
	Name   string
	People []*Person `linq:"many_to_many"`
}

type ProjectedPerson struct {
	Id   int
	Name string
}

type GroupByName struct {
	Name string
}

// for one to one testing
type OneToOnePerson struct {
	Id   int
	Name string
	Job  *OneToOneJob `linq:"one_to_one"`
}

type OneToOneJob struct {
	Id   int
	Name string
}

//for one to many testing

type OneToManyPerson struct {
	Id   int
	Name string
	Jobs []*OneToManyJob `linq:"one_to_many"`
}

type OneToManyJob struct {
	Id   int
	Name string
}

// for query testing

type Score struct {
	Id         int
	PlayerName string
	Value      int
}

type ProjectedScore struct {
	Id    int
	Value int
}

type GroupByPlayerName struct {
	PlayerName string
}
