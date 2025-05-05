package db

import (
	"log/slog"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gim/config"
	"gim/pkg/util"
)

var (
	DB        *gorm.DB
	RedisCli  *redis.Client
	RedisUtil *util.RedisUtil
)

func init() {
	InitMysql(config.Config.MySQL)
	InitRedis(config.Config.RedisHost, config.Config.RedisPassword)
}

// InitMysql 初始化MySQL
func InitMysql(dsn string) {
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
	DB = db
}

// InitRedis 初始化Redis
func InitRedis(addr, password string) {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       0,
		Password: password,
	})

	_, err := RedisCli.Ping().Result()
	if err != nil {
		slog.Error("redis ping error", "error", err)
		panic(err)
	}

	RedisUtil = util.NewRedisUtil(RedisCli)
}
