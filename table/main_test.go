package table

import (
	"fmt"
	"testing"
)

func assert(cond bool) {
    if !cond {
        panic(fmt.Sprintf("assertion failed"))
    }
}

func TestSimple(t *testing.T) {
    var st storage = &SimpleRamStorage{}
    int_field := IntField{"smth"}
    table := New(
    	"Main",
        st,
        &int_field,
    )
    err := table.Insert([]interface{}{int32(10000)})
    if err != nil {
        t.Error(err)
    }
    data, err := table.Select()
    if err != nil {
        t.Error(err)
    }
    assert((*data.Next())[0] == int32(10000))
    err = table.Insert([]interface{}{int32(20000)})
    if err != nil {
        t.Error(err)
    }
    err = table.Insert([]interface{}{int32(30000)})
    if err != nil {
        t.Error(err)
    }
    data, err = table.Select(&filter{
    	field: 0,
    	op: func(el interface{}) bool {
    		return el.(int32) >= int32(20000)
    	},
    })
    if err != nil {
        t.Error(err)
    }
    assert((*data.Next())[0] == int32(20000))
    assert((*data.Next())[0] == int32(30000))
    count, err := table.Delete(&filter{
    	field: 0,
    	op: func(el interface{}) bool {
    		return el.(int32) >= int32(20000)
    	},
    })
    if err != nil {
        t.Error(err)
    }
    assert(count == 2)
    data, err = table.Select()
    if err != nil {
        t.Error(err)
    }
    assert((*data.Next())[0] == int32(10000))
    assert(data.Next() == nil)
}
