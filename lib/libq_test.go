package lib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOneToOne(t *testing.T) {
	storage := NewStorage()

	job := &OneToOneJob{
		Name: "Test Job",
	}

	person := &OneToOnePerson{
		Name: "Test Person",
		Job:  job,
	}

	p, _ := storage.Store(person)
	storedPerson := p.(*OneToOnePerson)

	lj := storage.Load(storedPerson.Job.Id, OneToOneJob{})
	loadedJob := lj.(*OneToOneJob)

	assert.Equal(t, loadedJob.Name, "Test Job")
	assert.Equal(t, loadedJob.Id, storedPerson.Job.Id)
}

func TestOneToMany(t *testing.T) {
	storage := NewStorage()

	jobs := make([]*OneToManyJob, 0, 5)
	for i := 0; i < 5; i++ {
		jobs = append(jobs, &OneToManyJob{
			Name: fmt.Sprintf("Test Job %d", i),
		})
	}

	person := &OneToManyPerson{
		Name: "Test Person",
		Jobs: jobs,
	}

	p, _ := storage.Store(person)
	storedPerson := p.(*OneToManyPerson)

	for _, job := range storedPerson.Jobs {
		lj := storage.Load(job.Id, OneToManyJob{})
		loadedJob := lj.(*OneToManyJob)

		assert.Equal(t, job.Id, loadedJob.Id)
	}

}

func TestManyToMany(t *testing.T) {
	storage := NewStorage()
	job := &Job{
		Name: "student",
	}

	jobs := []*Job{job}

	person := Person{
		Name:    "Test",
		Surname: "Test",
		Jobs:    jobs,
	}

	personSec := Person{
		Name:    "Second",
		Surname: "Test",
		Jobs:    jobs,
	}

	job.People = []*Person{&person, &personSec}

	lp, _ := storage.Store(&person)

	storedPerson := lp.(*Person)

	storedJob := storedPerson.Jobs[0]
	assert.Equal(t, job.Id, 1)

	assert.Equal(t, storedJob.People[1].Name, "Second")
	sp := storage.Load(storedJob.People[1].Id, Person{})

	secondPerson := sp.(*Person)
	assert.Equal(t, secondPerson.Id, storedJob.People[1].Id)

}

func TestQueryCmd_Where(t *testing.T) {
	storage := NewStorage()
	for i := 0; i < 5; i++ {
		score := &Score{
			PlayerName: fmt.Sprintf("player %d", i),
			Value:      i,
		}

		storage.Store(score)
	}

	pipe := NewQueryCmd()

	out := pipe.Where(storage, Score{}, func(elem interface{}) bool {
		p := elem.(*Score)
		return p.Value%2 == 0
	}).Result()

	assert.Equal(t, len(out), 3)

}

func TestQueryCmd_Project(t *testing.T) {
	storage := NewStorage()
	for i := 0; i < 5; i++ {
		score := &Score{
			PlayerName: fmt.Sprintf("player %d", i),
			Value:      i,
		}

		storage.Store(score)
	}

	pipe := NewQueryCmd()

	out := pipe.Project(storage, Score{}, ProjectedScore{}).Result()
	assert.Equal(t, len(out), 5)

	firstElemCasted := out[0]

	assert.IsType(t, ProjectedScore{}, firstElemCasted)
}

func TestQueryCmd_GroupBy(t *testing.T) {
	storage := NewStorage()
	for i := 0; i < 5; i++ {
		score := &Score{
			PlayerName: "player",
			Value:      i,
		}

		storage.Store(score)
	}

	pipe := NewQueryCmd()

	out := pipe.GroupBy(storage, Score{}, GroupByPlayerName{}, 0, func(acc any, elem interface{}) any {
		ps := elem.(*Score)
		cur := acc.(int)
		return cur + ps.Value
	}).Result()

	assert.Equal(t, len(out), 1)

	result := out[0].(GroupByResult)
	assert.Equal(t, result.Result, 10) // 0 + 1 + 2 + 3 + 4
}

func TestComplexQuery(t *testing.T) {
	storage := NewStorage()
	for i := 0; i < 5; i++ {
		score := &Score{
			PlayerName: "player",
			Value:      i,
		}

		storage.Store(score)
	}

	pipe := NewQueryCmd()

	out := pipe.
		Where(storage, Score{}, func(elem interface{}) bool {
			p := elem.(*Score)
			return p.Value%2 == 0
		}).
		Project(storage, Score{}, ProjectedScore{}).
		GroupBy(storage, Score{}, GroupByPlayerName{}, 0, func(acc any, elem interface{}) any {
			ps := elem.(ProjectedScore)
			cur := acc.(int)
			return cur + ps.Value
		}).Result()

	assert.Equal(t, len(out), 1)

	result := out[0].(GroupByResult)
	assert.Equal(t, result.Result, 6) // 0 + 2 + 4
}
