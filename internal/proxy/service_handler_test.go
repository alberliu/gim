package proxy

import (
	"context"
	"testing"

	"gim/pkg/gen/proto/businesspb"
	"gim/pkg/local"
)

func TestAuth(t *testing.T) {
	_, err := local.GetMessageExtClient().SendFriendMessage(context.Background(), &businesspb.SendFriendMessageRequest{
		UserId:  0,
		Content: nil,
	})
	if err != nil {
		t.Error(err)
	}
}
