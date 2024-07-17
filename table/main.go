package table

import "github.com/Visoff/DataBase/iterator"

type Table struct {
    Name []byte
    header header
    data storage
}

type header struct {
    fields []field
    row_size int
}

func (h *header)Values(row [][]byte) ([]interface{}, error) {
    res := make([]interface{}, 0)
    for i, field := range h.fields {
        val, err := field.Value(row[i])
        if err != nil {
            return nil, err
        }
        res = append(res, val)
    }
    return res, nil
}

type field interface {
    Name() []byte
    Size() int
    Value([]byte) (interface{}, error)
    Store(interface{}) ([]byte, error)
}

type storage interface {
    Get() (iterator.Iterator[[][]byte], error)
    Insert([][]byte) error
}

func New(name string, storage storage, fields ...field) *Table {
    row_size := 0
    for _, field := range fields {
        row_size += field.Size()
    }
    return &Table{
        Name: []byte(name),
        header: header{
            fields: fields,
            row_size: row_size,
        },
        data: storage,
    }
}

func (t *Table) Insert(it []interface{}) error {
    data := make([][]byte, len(t.header.fields))
    for i, field := range t.header.fields {
        val, err := field.Store(it[i])
        if err != nil {
            return err
        }
        data[i] = val
    }
    return t.data.Insert(data)
}

type filter struct {
    field int
    op func(interface{}) bool
}

func (t *Table) run_filter(f *filter, row [][]byte) bool {
    field := t.header.fields[f.field]
    value, err := field.Value(row[f.field])
    if err != nil {
        return false
    }
    return f.op(value)
}

func (t *Table) Select(filters ...*filter) ([][]interface{}, error) {
    res := make([][]interface{}, 0)
    rows, err := t.data.Get()
    if err != nil {
        return nil, err
    }
    for row := rows.Next(); row != nil; row = rows.Next() {
        accept := true
        row_data := *row
        for _, filter := range filters {
            if !t.run_filter(filter, row_data) {
                accept = false
                break
            }
        }
        if accept {
            val, err := t.header.Values(row_data)
            if err != nil {
                return nil, err
            }
            res = append(res, val)
        }
    }
    
    return res, nil
}
