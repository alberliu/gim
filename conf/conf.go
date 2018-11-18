package conf

import "os"

// MySQL mysql配置
var MySQL = "root:Liu123456@tcp(localhost:3306)/im?charset=utf8&parseTime=true"

func init() {
	env := os.Getenv("im_env")
	if env == "dev" {
		initDevelopConf()
	}

	if env == "pro" {
		initProductConf()
	}
}

func initDevelopConf() {

}

func initProductConf() {

}
