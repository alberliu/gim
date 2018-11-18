package db

import (
	"goim/public/session"

	"goim/conf"

	_ "github.com/go-sql-driver/mysql"
)

var Factoty *session.SessionFactory

func init() {
	var err error
	Factoty, err = session.NewSessionFactory("mysql", conf.MySQL)
	if err != nil {
		panic(err)
	}
}
