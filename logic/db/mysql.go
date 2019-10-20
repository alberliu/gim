package db

import (
	"database/sql"
	"goim/conf"

	_ "github.com/go-sql-driver/mysql"
)

var DBCli *sql.DB

func init() {
	var err error
	DBCli, err = sql.Open("mysql", conf.MySQL)
	if err != nil {
		panic(err)
	}
}
