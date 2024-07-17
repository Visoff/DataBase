package main

import (
    "github.com/Visoff/DataBase/btree"
)

func main() {
    t := btree.New[int, string](3)
    _ = t
}
