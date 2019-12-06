package conf

import "os"

// conn和logic公用配置
var (
	MySQL   = "root:liu123456@tcp(localhost:3306)/gim?charset=utf8&parseTime=true"
	NSQIP   = "127.0.0.1:4150"
	RedisIP = "127.0.0.1:6379"
)

// conn配置
var (
	ConnTCPListenAddr = ":8080"
	ConnRPCListenAddr = ":60000"
	LocalAddr         = "127.0.0.1:60000"
	LogicRPCAddrs     = "addrs:///127.0.0.1:50000"
)

// logic配置
var (
	LogicRPCIntListenAddr       = ":50000"
	LogicClientRPCExtListenAddr = ":50001"
	LogicServerRPCExtListenAddr = ":50002"
	ConnRPCAddrs                = "addrs:///127.0.0.1:60000"
)

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
