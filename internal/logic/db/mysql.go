package db

import (
	"database/sql"
	"gim/config"

	_ "github.com/go-sql-driver/mysql"
)

var DBCli *sql.DB

func init() {
	var err error
	DBCli, err = sql.Open("mysql", config.LogicConf.MySQL)
	if err != nil {
		panic(err)
	}
}
