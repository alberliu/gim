package dao

import (
	"fmt"
	"gim/pkg/db"
	"gim/pkg/logger"
)

func init() {
	fmt.Println("init db")
	logger.Target = logger.Console
	logger.Init()
	db.InitByTest()
}
