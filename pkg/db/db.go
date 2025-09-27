package db

import (
	"log/slog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gim/config"
	"gim/pkg/uredis"
)

var (
	DB       *gorm.DB
	RedisCli *uredis.Client
)

func init() {
	DB = newDB(config.Config.MySQL)
	RedisCli = uredis.NewClient(config.Config.RedisHost, config.Config.RedisPassword)
}

func newDB(dsn string) *gorm.DB {
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		slog.Error("open db error", "error", err, slog.String("dsn", dsn))
		panic(err)
	}

	if config.ENV == config.EnvLocal {
		db = db.Debug()
	}
	return db
}
