package utils

type Empty struct{}

var empty Empty

type Set struct {
	m map[interface{}]Empty
}

func NewSet() *Set {
	return &Set{
		m: make(map[interface{}]Empty),
	}
}

//添加元素
func (s *Set) Add(val interface{}) {
	s.m[val] = empty
}

//删除元素
func (s *Set) Remove(val interface{}) {
	delete(s.m, val)
}

//获取长度
func (s *Set) Len() interface{} {
	return len(s.m)
}

func (s *Set) Contains(val interface{}) bool {
	_, ok := s.m[val]
	return ok
}

//清空set
func (s *Set) Clear() {
	s.m = make(map[interface{}]Empty)
}

func (s *Set) Array() (ret []interface{}) {
	for k := range s.m {
		ret = append(ret, k)
	}
	return
}
