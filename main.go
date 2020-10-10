package main

import (
	"db"
	"fmt"
)

func main() {
	d := new(db.DB)
	d.Init()
	ok, err := d.Insert(&db.Pair{Key: 1, Val: []byte("123")})
	if !ok || err != nil {
		fmt.Println("insert failed.")
	}
	pair := d.Search(1)
	if pair == nil {
		fmt.Println("search failed.")
	} else {
		fmt.Printf("search found: %v\n", pair)
	}
}
