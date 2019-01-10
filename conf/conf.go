package conf

import "os"

// MySQL mysql配置
var MySQL = "ximi:ximi.2018@tcp(localhost:3306)/im_bak?charset=utf8&parseTime=true"

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
