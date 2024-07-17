package btree

import (
	"cmp"
	"errors"
	"fmt"
	"strings"
)

var KeyNotFoundError = errors.New("key not found")

type Data[K cmp.Ordered, V any] struct {
    key K
    value V
}

type node[K cmp.Ordered, V any] struct {
    keys []Data[K, V]
    children []*node[K, V]
}

type BTree[K cmp.Ordered, V any] struct{
    root *node[K, V]
    min_angle int
}

func New[K cmp.Ordered, V any](min_angle int) *BTree[K, V] {
    return &BTree[K, V]{
        root: &node[K, V]{
            keys: make([]Data[K, V], 0),
            children: make([]*node[K, V], 0),
        },
        min_angle: min_angle,
    }
}

func (t *BTree[K, V]) Find(key K) (*V, error) {
    n, _ := t.root.find(key)
    if n == nil {
        return nil, KeyNotFoundError
    }
    return n.get_value(key)
}

func (n *node[K, V]) leaf() bool {
    return len(n.children) == 0
}

func (n *node[K, V]) find(key K) (*node[K, V], int) {
    i := 0
    for i < len(n.keys) && key > n.keys[i].key {
        i++
    }
    if i < len(n.keys) && key == n.keys[i].key {
        return n, i
    }
    if n.leaf() {
        return nil, 0
    }
    return n.children[i].find(key)
}

func (n *node[K, V]) get_value(key K) (*V, error) {
    for _, v := range n.keys {
        if key == v.key {
            return &v.value, nil
        }
    }
    return nil, KeyNotFoundError
}

func (t *BTree[K, V]) Insert(key K, value V) {
    t.split(t.root, 0, true)
    ptr := t.root
    i := 0
    for i < len(ptr.keys) && key > ptr.keys[i].key {
        i++
    }
    for !ptr.leaf() {
        t.split(ptr, i, false)
        for i < len(ptr.keys) && key > ptr.keys[i].key {
            i++
        }
        ptr = ptr.children[i]
    }
    i = 0
    for i < len(ptr.keys) && key > ptr.keys[i].key {
        i++
    }
    ptr.keys = append(ptr.keys[:i], append([]Data[K, V]{{key, value}}, ptr.keys[i:]...)...)
}

func (t *BTree[K, V]) split(p *node[K, V], i int, root bool) {
    if root {
        if len(p.keys) < t.min_angle*2 - 1 {return}
        t.root = &node[K, V]{
            keys: make([]Data[K, V], 0),
            children: []*node[K, V]{t.root},
        }
        p = t.root
        i = 0
    }
    n := p.children[i]
    nk := len(n.keys)
    if nk < t.min_angle*2 - 1 {return}
    mid := t.min_angle-1
    left := node[K, V]{
        keys: n.keys[:mid],
    }
    right := node[K, V]{
        keys: n.keys[mid+1:],
    }
    p.keys = append(p.keys[:i], append([]Data[K, V]{n.keys[mid]}, p.keys[i:]...)...)
    p.children = append(p.children[:i], append([]*node[K, V]{&left, &right}, p.children[i+1:]...)...)
}

func display[K cmp.Ordered, V any](n *node[K, V]) string {
    if n == nil {return ""}
    res := make([]string, 0)
    for _, data := range n.keys {
        res = append(res, fmt.Sprint(data.key))
    }
    children := make([]string, 0)
    for _, chld := range n.children {
        children = append(children, display(chld))
    }
    return "("+strings.Join(res, ", ")+") -> {"+strings.Join(children, ", ")+"}"
}
