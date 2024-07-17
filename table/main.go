package table

import (
	"github.com/Visoff/DataBase/iterator"
)

type Table struct {
    Name []byte
    header header
    data storage
}

type header struct {
    fields []field
    row_size int
    primary_keys []int
}

func (h *header)ToValues(row [][]byte) ([]interface{}, error) {
    res := make([]interface{}, len(h.fields))
    for i, field := range h.fields {
        val, err := field.Value(row[i])
        if err != nil {
            return nil, err
        }
        res[i] = val
    }
    return res, nil
}

func (h *header)FromValues(values []interface{}) ([][]byte, error) {
    res := make([][]byte, len(h.fields))
    for i, field := range h.fields {
        val, err := field.Store(values[i])
        if err != nil {
            return nil, err
        }
        res[i] = val
    }
    return res, nil
}

type field interface {
    Name() []byte
    Size() int
    Value([]byte) (interface{}, error)
    Store(interface{}) ([]byte, error)

    Hash([]byte) (uint64, error)
}

type storage interface {
    Init(*header)

    Get([]*filter, []*filter) (iterator.Iterator[[]interface{}], error)
    Insert([]interface{}) error
    Delete([]*filter, []*filter) (int, error)
}

func New(name string, st storage, fields ...field) *Table {
    row_size := 0
    for _, field := range fields {
        row_size += field.Size()
    }
    h := header {
        fields: fields,
        row_size: row_size,
    }
    st.Init(&h)
    return &Table{
        Name: []byte(name),
        header: h,
        data: st,
    }
}

func (t *Table) Insert(data []interface{}) error {
    return t.data.Insert(data)
}

type filter struct {
    field int
    op func(interface{}) bool
}

func (t *Table) split_filters(filters []*filter) ([]*filter, []*filter) {
    pk_filtes := make([]*filter, 0)
    other_filtes := make([]*filter, 0)
    for _, f := range filters {
        contains := false
        for _, pk := range t.header.primary_keys {
            if pk == f.field {
                contains = true
                break
            }
        }
        if contains {
            pk_filtes = append(pk_filtes, f)
        } else {
            other_filtes = append(other_filtes, f)
        }
    }
    return pk_filtes, other_filtes
}

func (t *Table) Select(filters ...*filter) (iterator.Iterator[[]interface{}], error) {
    pk_filtes, other_filtes := t.split_filters(filters)
    return t.data.Get(pk_filtes, other_filtes)
}

func (t *Table) Delete(filters ...*filter) (int, error) {
    pk_filtes, other_filtes := t.split_filters(filters)
    return t.data.Delete(pk_filtes, other_filtes)
}
