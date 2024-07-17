package iterator

import "fmt"

type node[T any] struct {
    value T
    next *node[T]
    prev *node[T]
}

type Iterator[T any] struct {
    start *node[T]
    end *node[T]
}

func New[T any]() Iterator[T] {
    return Iterator[T]{
        start: nil,
        end: nil,
    }
}

func FromArray[T any](arr []T) Iterator[T] {
    it := New[T]()
    for _, v := range arr {
        it.PushLeft(v)
    }
    return it
}

func (i *Iterator[T])ToArray() []T {
    res := make([]T, 0)
    for n := i.start; n != nil; n = n.next {
        res = append(res, n.value)
    }
    i.start = nil
    i.end = nil
    return res
}

func (n *node[T])display() string {
    if n == nil {return "nil"}
    return fmt.Sprint(n.value)+" -> "+n.next.display()
}

func (i *Iterator[T])Next() *T {
    if i.start == nil {return nil}
    res := i.start.value
    i.start = i.start.next
    if i.start == nil {
        i.end = nil
        return &res
    }
    i.start.prev = nil
    return &res
}

func (i *Iterator[T])Peek() *T {
    if i.start == nil {return nil}
    return &i.start.value
}

func (i *Iterator[T])PushLeft(val T) {
    n := &node[T]{
        value: val,
        next: i.start,
        prev: nil,
    }
    if i.start != nil {
        i.start.prev = n
        i.start = n
    }
    if i.start == nil || i.end == nil {
        i.start = n
        i.end = n
    }
    i.start = n
}

func (i *Iterator[T])PushRight(val T) {
    n := &node[T]{
        value: val,
        next: nil,
        prev: i.end,
    }
    if i.end != nil {
        i.end.next = n
        i.end = n
    }
    if i.start == nil || i.end == nil {
        i.start = n
        i.end = n
    }
    i.end = n
}
