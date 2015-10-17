// go-rdb
package core

import "reflect"

// Join
type JoinStream struct {
	Input1   Stream   `json:"input1"`
	Input2   Stream   `json:"input2"`
	Attr1    string   `json:"attr1"`
	Attr2    string   `json:"attr2"`
	Selector Operator `json:"selector"`

	index        int
	tuples       []*Tuple
	currentTuple *Tuple
	currentKind  reflect.Kind
	targetKind   reflect.Kind
}

func (s *JoinStream) Next() (result *Tuple, err error) {
	if len(s.tuples) <= s.index {
		s.index = 0
		s.currentTuple = nil
	}
	if s.currentTuple == nil {
		if s.Input1.HasNext() {
			s.currentTuple, err = s.Input1.Next()
		}
		if s.currentTuple == nil {
			return
		}
		s.currentKind = s.currentTuple.Schema.GetKind(s.Attr1)
	}
	targetTuple := s.tuples[s.index]
	if s.targetKind == 0 {
		s.targetKind = targetTuple.Schema.GetKind(s.Attr2)
	}
	s.index++
	res, err := s.Selector(s.currentKind, s.currentTuple.Get(s.Attr1), s.targetKind, targetTuple.Get(s.Attr2))
	if err != nil {
		return
	}
	if res {
		result = NewTuple()
		for i, attr := range s.currentTuple.Attrs {
			result.Set(attr, s.currentTuple.Data[i])
		}
		for i, attr := range targetTuple.Attrs {
			result.Set(attr, targetTuple.Data[i])
		}
		return
	}
	if s.HasNext() {
		return s.Next()
	}
	return nil, nil
}
func (s *JoinStream) HasNext() bool {
	if s.tuples == nil {
		s.tuples = make([]*Tuple, 0, TupleCapacity)
		for s.Input2.HasNext() {
			next, err := s.Input2.Next()
			if err != nil {
				continue
			}
			s.tuples = append(s.tuples, next)
		}
	}
	if len(s.tuples) > s.index {
		return true
	}
	return s.Input1.HasNext()
}
func (s *JoinStream) Init(n *Node) error {
	if err := s.Input1.Init(n); err != nil {
		return err
	}
	return s.Input2.Init(n)
}
func (s *JoinStream) Close() {
	s.Input1.Close()
	s.Input2.Close()
}
