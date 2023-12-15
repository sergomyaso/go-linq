package lib

import (
	"fmt"
	"reflect"
)

const (
	OneToOne   = "one_to_one"
	OneToMany  = "one_to_many"
	ManyToMany = "many_to_many"
)

type Table struct {
	currentIndex int
	store        map[int]interface{}
}

func NewTable() *Table {
	return &Table{
		currentIndex: 1,
		store:        make(map[int]interface{}),
	}
}

func (t *Table) Add(v interface{}) int {
	t.store[t.currentIndex] = v
	t.currentIndex++

	return t.currentIndex - 1
}

func (t *Table) Set(id int, v interface{}) {
	t.store[id] = v
}

func (t *Table) Delete(id int) {
	delete(t.store, id)
}

type Storage struct {
	tables map[string]*Table // TypeName -> Table
}

func NewStorage() *Storage {
	return &Storage{
		tables: make(map[string]*Table),
	}
}

func (s *Storage) Store(value interface{}) (interface{}, error) {
	entityId := int(reflect.ValueOf(value).Elem().FieldByName("Id").Int())

	value, err := s.storeRelations(value)
	if err != nil {
		return nil, err
	}

	valueType := reflect.TypeOf(value).Elem().Name()

	// получаем нужную таблицу по типу
	table, ok := s.tables[valueType]
	if !ok {
		table = NewTable()
		s.tables[valueType] = table
	}

	if entityId != 0 {
		table.Set(entityId, value)
		return value, nil
	}

	entityId = table.Add(value)

	// сохраняем новый id
	reflect.ValueOf(value).Elem().FieldByName("Id").SetInt(int64(entityId))

	return value, nil
}

func (s *Storage) Load(id int, target any) interface{} {
	typeName := reflect.TypeOf(target).Name()

	table, ok := s.tables[typeName]
	if !ok {
		return nil
	}

	return table.store[id]
}

func (s *Storage) Delete(id int, target any) {
	typeName := reflect.TypeOf(target).Name()

	table, ok := s.tables[typeName]
	if !ok {
		return
	}

	table.Delete(id)
}

func (s *Storage) Where(target any, filter func(elem interface{}) bool) []interface{} {
	typeName := reflect.TypeOf(target).Name()
	table, ok := s.tables[typeName]
	if !ok {
		return nil
	}

	output := make([]interface{}, 0)
	for _, value := range table.store {
		if filter(value) {
			output = append(output, value)
		}
	}

	return output
}

// Для сохранения зависимых объектов
func (s *Storage) storeRelations(value interface{}) (interface{}, error) {
	valueType := reflect.TypeOf(value).Elem()

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)

		tag := field.Tag.Get("linq")

		switch tag {
		case OneToOne:
			stored, err := s.Store(reflect.ValueOf(value).Elem().Field(i).Interface())
			if err != nil {
				return nil, err
			}

			reflect.ValueOf(value).Elem().Field(i).Set(reflect.ValueOf(stored))
		case OneToMany:
			sl := reflect.ValueOf(value).Elem().Field(i)
			if sl.Kind() != reflect.Slice {
				return nil, fmt.Errorf("one_to_many related object is not a slice")
			}

			for j := 0; j < sl.Len(); j++ {
				fmt.Println(sl.Index(j).Elem().Type())
				stored, err := s.Store(sl.Index(j).Interface())
				if err != nil {
					return nil, err
				}

				reflect.ValueOf(value).Elem().Field(i).Index(j).Set(reflect.ValueOf(stored))
			}
		}
	}

	return value, nil
}
