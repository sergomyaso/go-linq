package lib

type Person struct {
	Id      int
	Name    string
	Surname string
	Jobs    []*Job `linq:"one_to_many"`
}

type Job struct {
	Id   int
	Name string
}
