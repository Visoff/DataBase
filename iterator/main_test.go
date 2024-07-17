package iterator

import (
	"testing"
)

func assert(b bool) {
    if !b {
        panic("assertion failed")
    }
}

func TestInit(t *testing.T) {
    it := New[int]()
    assert(it.start == nil)
    assert(it.end == nil)
}

func TestMain(t *testing.T) {
    it := New[int]()
    it.PushRight(1)
    it.PushRight(2)
    assert(*it.Next() == 1)
    assert(*it.Next() == 2)
    it.PushLeft(1)
    it.PushLeft(2)
    assert(*it.Next() == 2)
    assert(*it.Next() == 1)
}
