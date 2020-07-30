package dao

import (
	"fmt"
	"gim/pkg/db"
)

func init() {
	fmt.Println("init db")
	db.InitByTest()
}
