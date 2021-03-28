package cache

import (
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

type queue struct{}

var Queue = new(queue)

func (queue) Publish(topic string, bytes []byte) error {
	_, err := db.RedisCli.Publish(topic, bytes).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}
