package lib

import (
	"goim/public/logger"

	"github.com/json-iterator/go"
)

func JsonMarshal(v interface{}) string {
	bytes, err := jsoniter.Marshal(v)
	if err != nil {
		logger.Sugar.Error("json序列化：", err)
	}
	return Bytes2str(bytes)
}
