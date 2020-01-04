package uid

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestLid(t *testing.T) {
	db, err := sql.Open("mysql", "root:Liu123456@tcp(localhost:3306)/im?charset=utf8")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	lid, err := NewLid(db, "test", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	i := 0
	for i < 100 {
		id := lid.Get()
		fmt.Println(id)
		i++
	}
}

func TestLid_Get(t *testing.T) {
	go getLid("one")
	go getLid("two")
	go getLid("three")
	select {}
}

func getLid(index string) {
	db, err := sql.Open("mysql", "root:Liu123456@tcp(localhost:3306)/im?charset=utf8")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	lid, err := NewLid(db, "test", 1000)
	if err != nil {
		fmt.Println(err)
		return
	}
	i := 0
	for i < 100 {
		id := lid.Get()
		fmt.Println(index, id)
		i++
	}
}

func BenchmarkLeafKey(b *testing.B) {
	db, err := sql.Open("mysql", "root:Liu123456@tcp(localhost:3306)/im?charset=utf8")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	lid, err := NewLid(db, "test", 1000)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < b.N; i++ {
		lid.Get()
	}
}
