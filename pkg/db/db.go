package db

import (
	"gim/config"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB       *gorm.DB
	RedisCli *redis.Client
)

// InitMysql 初始化MySQL
func InitMysql(dataSource string) {
	var err error
	DB, err = gorm.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	DB.SingularTable(true)
	DB.LogMode(true)
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
		panic(err)
	}
}

// InitByTest 初始化数据库配置，仅用在单元测试
func InitByTest() {
	InitMysql(config.Logic.MySQL)
	InitRedis(config.Logic.RedisIP, config.Logic.RedisPassword)
}
