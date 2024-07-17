package table

import (
	"github.com/Visoff/DataBase/iterator"
)

type SimpleRamStorage struct {
	data [][][]byte
    header *header
}

func (s *SimpleRamStorage) Init(header *header) {
    s.header = header
}

func (s *SimpleRamStorage) Insert(row []interface{}) error {
    data, err := s.header.FromValues(row)
    if err != nil {
        return err
    }
	s.data = append(s.data, data)
	return nil
}

func (s *SimpleRamStorage) run_filter(f *filter, row [][]byte) bool {
    field := s.header.fields[f.field]
    value, err := field.Value(row[f.field])
    if err != nil {
        return false
    }
    return f.op(value)
}

func (s *SimpleRamStorage) Get(pk_filters []*filter, other_filters []*filter) (iterator.Iterator[[]interface{}], error) {
    res := iterator.New[[]interface{}]()
    filters := append(pk_filters, other_filters...)
    for _, row := range s.data {
        accepted := true
        for _, f := range filters {
            if !s.run_filter(f, row) {
                accepted = false
                break
            }
        }
        if accepted {
            values, err := s.header.ToValues(row)
            if err != nil {
                return res, err
            }
            res.PushRight(values)
        }
    }
    return res, nil
}

func (s *SimpleRamStorage) Delete(pk_filters []*filter, other_filters []*filter) (int, error) {
    res := 0
    filters := append(pk_filters, other_filters...)
    for i := len(s.data) - 1; i >= 0; i-- {
        accepted := true
        for _, f := range filters {
            if !s.run_filter(f, s.data[i]) {
                accepted = false
                break
            }
        }
        if accepted {
            res++
            s.data = append(s.data[:i], s.data[i+1:]...)
        }
    }
    return res, nil
}

type IntField struct {
    name string
}

func (f *IntField) Name() []byte {
    return []byte(f.name)
}

func (f *IntField) Size() int {
    return 4
}

func (f *IntField) Value(data []byte) (interface{}, error) {
    return int32(data[0]) | int32(data[1])<<8 | int32(data[2])<<16 | int32(data[3])<<24, nil
}

func (f *IntField) Store(value interface{}) ([]byte, error) {
    return []byte{
        byte(value.(int32)),
        byte(value.(int32) >> 8),
        byte(value.(int32) >> 16),
        byte(value.(int32) >> 24),
    }, nil
}

func (f *IntField) Hash(data []byte) (uint64, error) {
    res, err := f.Value(data)
    if err != nil {
        return 0, err
    }
    return uint64(res.(int32)), nil
}
