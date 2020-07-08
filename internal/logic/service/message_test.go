package service

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestMessageService_Add(t *testing.T) {

}

func TestMessageService_ListByUserIdAndSequence(t *testing.T) {

}

func TestJson(t *testing.T) {
	var st = struct {
		Nickname string `json:"nickname"`
	}{}

	json := `{
	"user_id":3,
	"nickname":"h",
	"sex":2,
	"avatar_url":"no",
	"extra":{"nickname":"hjkladsjfkl"}
}`
	jsoniter.Get([]byte(json), "extra").ToVal(&st)
	fmt.Println(st)
}
