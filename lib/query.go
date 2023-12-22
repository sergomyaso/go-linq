package lib

import "reflect"

type Start struct{}

type QueryCmd struct {
	cmds  []<-chan []interface{}
	start chan Start
	curQ  int
}

func NewQueryCmd() *QueryCmd {
	return &QueryCmd{
		cmds:  make([]<-chan []interface{}, 0),
		start: make(chan Start),
	}
}

type GroupByResult struct {
	GroupBy any
	Result  any
}

func (q *QueryCmd) Result() []interface{} {
	q.start <- Start{}
	output := <-q.cmds[q.curQ-1]
	return output
}

func (q *QueryCmd) Where(s *Storage, target any, filter func(elem interface{}) bool) *QueryCmd {
	res := make(chan []interface{})
	go func(cur int) {
		if cur == 0 {
			select {
			case <-q.start:
				break
			}
		} else {
			select {
			case input := <-q.cmds[cur-1]:
				output := make([]interface{}, 0)
				for _, value := range input {
					if filter(value) {
						output = append(output, value)
					}
				}

				res <- output
				return
			}
		}

		typeName := reflect.TypeOf(target).Name()
		table, ok := s.tables[typeName]
		if !ok {
			res <- nil
		}

		output := make([]interface{}, 0)
		for _, value := range table.store {
			if filter(value) {
				output = append(output, value)
			}
		}

		res <- output
	}(q.curQ)

	q.cmds = append(q.cmds, res)
	q.curQ++

	return q
}

func (q *QueryCmd) GroupBy(s *Storage, target any, groupBy any, acc any, op func(acc any, elem interface{}) any) *QueryCmd {
	res := make(chan []interface{})
	go func(cur int) {
		if cur == 0 {
			select {
			case <-q.start:
				break
			}
		} else {
			select {
			case input := <-q.cmds[cur-1]:
				// обработка пайплайна
				groups := GetGroupsFromPipe(input, groupBy)
				output := make([]interface{}, 0, len(groups))

				for key, group := range groups {
					groupAcc := acc
					for _, elem := range group {
						groupAcc = op(groupAcc, elem)
					}

					output = append(output, GroupByResult{
						GroupBy: key,
						Result:  groupAcc,
					})

				}

				res <- output
				return
			}
		}

		// обработка пайплайна
		groups := GetGroupsFromStorage(s, target, groupBy)
		output := make([]interface{}, 0, len(groups))
		for key, group := range groups {

			for _, elem := range group {
				acc = op(acc, elem)
			}

			output = append(output, GroupByResult{
				GroupBy: key,
				Result:  acc,
			})

		}

		res <- output
	}(q.curQ)

	q.cmds = append(q.cmds, res)
	q.curQ++

	return q
}

func (q *QueryCmd) Project(s *Storage, target any, projectStructure any) *QueryCmd {
	res := make(chan []interface{})
	go func(cur int) {
		if cur == 0 {
			select {
			case <-q.start:
				break
			}
		} else {
			select {
			case input := <-q.cmds[cur-1]:
				// обработка пайплайна
				output := make([]interface{}, 0)
				for _, value := range input {
					output = append(output, MapStructs(value, projectStructure))
				}

				res <- output
				return
			}
		}

		// обработка пайплайна
		typeName := reflect.TypeOf(target).Name()
		table, ok := s.tables[typeName]
		if !ok {
			res <- nil
		}

		output := make([]interface{}, 0)
		for _, value := range table.store {
			output = append(output, MapStructs(value, projectStructure))
		}

		res <- output
	}(q.curQ)

	q.cmds = append(q.cmds, res)
	q.curQ++

	return q
}

func MapStructs[S any, D any](src S, dest D) D {
	srcVal := reflect.Indirect(reflect.ValueOf(src))
	destVal := reflect.Indirect(reflect.ValueOf(dest))

	destPtr := reflect.New(destVal.Type())
	destPtr.Elem().Set(destVal)

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcFieldName := srcVal.Type().Field(i).Name

		destField := destPtr.Elem().FieldByName(srcFieldName)
		if destField.IsValid() && destField.Type() == srcField.Type() {
			destField.Set(srcField)
		}
	}

	return destPtr.Elem().Interface().(D)
}

func GetGroupsFromPipe(input []any, groupBy any) map[any][]any {
	groups := make(map[any][]any)
	for _, value := range input {
		projected := MapStructs(value, groupBy)
		groups[projected] = append(groups[projected], value)
	}

	return groups
}

func GetGroupsFromStorage(s *Storage, target any, groupBy any) map[any][]any {
	typeName := reflect.TypeOf(target).Name()
	table, ok := s.tables[typeName]
	if !ok {
		return nil
	}

	groups := make(map[any][]any)
	for _, value := range table.store {
		projected := MapStructs(value, groupBy)
		groups[projected] = append(groups[projected], value)
	}

	return groups
}
