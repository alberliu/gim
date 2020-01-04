package util

import (
	"encoding/json"
	"gim/pkg/logger"

	"go.uber.org/zap"
)

func JsonMarshal(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		logger.Logger.Error("json序列化：", zap.Error(err))
	}
	return Bytes2str(bytes)
}
