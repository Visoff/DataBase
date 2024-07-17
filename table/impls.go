package table

import "github.com/Visoff/DataBase/iterator"

type RamStorage struct {
	data [][][]byte
}

func (s *RamStorage) Set(i int, data [][]byte) error {
	s.data[i] = data
	return nil
}

func (s *RamStorage) Insert(data [][]byte) error {
	s.data = append(s.data, data)
	return nil
}

func (s *RamStorage) Delete(i int) error {
	s.data = append(s.data[:i], s.data[i+1:]...)
	return nil
}

func (s *RamStorage) Get() (iterator.Iterator[[][]byte], error) {
    return iterator.FromArray(s.data), nil
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
