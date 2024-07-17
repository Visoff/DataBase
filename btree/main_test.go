package btree

import (
	"testing"
)

func assert(b bool) {
    if !b {panic("assert panic")}
}

func TestMain(t *testing.T) {
    _ = New[int, string](3)
}

func TestSplit(t *testing.T) {
    bt := New[int, string](2)
    bt.root = &node[int, string]{
        keys: []Data[int, string]{{1, "a"}, {2, "b"}, {3, "c"}},
        children: make([]*node[int, string], 0),
    }
    bt.split(bt.root, 0, true)
    assert(bt.root.keys[0].key == 2)
    assert(bt.root.children[0].keys[0].key == 1)
    assert(bt.root.children[1].keys[0].key == 3)
}

func TestInsert(t *testing.T) {
    bt := New[int, string](2)
    bt.Insert(1, "1")
    bt.Insert(2, "2")
    bt.Insert(3, "3")
    bt.Insert(4, "4")
    bt.Insert(5, "5")
    bt.Insert(6, "6")
    bt.Insert(7, "7")
    bt.Insert(8, "8")
    assert(display(bt.root) == "(2, 4, 6) -> {(1) -> {}, (3) -> {}, (5) -> {}, (7, 8) -> {}}")
}

func TestFind(t *testing.T) {
    bt := New[int, string](2)
    bt.Insert(1, "1")
    bt.Insert(2, "2")
    bt.Insert(3, "3")
    bt.Insert(4, "4")
    bt.Insert(5, "5")
    bt.Insert(6, "6")
    bt.Insert(7, "7")
    bt.Insert(8, "8")
    v, err := bt.Find(4)
    if err != nil {
        t.Error(err)
    }
    assert(*v == "4")
}
