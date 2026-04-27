package proxy

import (
	"context"
	"testing"

	"gim/pkg/local"
	"gim/pkg/protocol/pb/businesspb"
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
