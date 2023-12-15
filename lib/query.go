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
