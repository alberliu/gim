package service

import (
	"context"
	"fmt"
	"gim/pkg/db"
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

func Test_messageService_Sync(t *testing.T) {
	db.InitByTest()
	resp, err := MessageService.Sync(context.TODO(), 6, 0)
	fmt.Println(err)
	fmt.Println(resp.HasMore)
	fmt.Println(len(resp.Messages))
}
