package service

import (
	"fmt"
	"gim/logic/model"
	"gim/public/pb"
	"testing"

	"github.com/json-iterator/go"
)

func TestMessageToJson(t *testing.T) {
	message := model.SendMessage{
		MessageBody: &pb.MessageBody{
			MessageContent: &pb.MessageContent{},
		},
	}
	bytes, err := jsoniter.Marshal(message)
	fmt.Println(err)
	fmt.Println(string(bytes))
}

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
