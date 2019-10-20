package util

import (
	"goim/public/logger"

	"go.uber.org/zap"

	"github.com/json-iterator/go"
)

func JsonMarshal(v interface{}) string {
	bytes, err := jsoniter.Marshal(v)
	if err != nil {
		logger.Logger.Error("json序列化：", zap.Error(err))
	}
	return Bytes2str(bytes)
}
