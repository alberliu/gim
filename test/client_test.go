package test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gim/pkg/gen/proto/businesspb"
	"gim/pkg/local"
	"gim/pkg/logger"
	"gim/pkg/md"
)

func TestClient(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: logger.ReplaceAttr,
	})))
	var header metadata.MD

	initData()

	network := NetworkWebsocket
	user1 := connect(network, 1, 11)
	connect(network, 1, 12)
	connect(network, 2, 2)
	connect(network, 3, 3)

	ctx := generateCtx(user1)
	messageClient := local.GetMessageExtClient()

	time.Sleep(2 * time.Second)
	fmt.Println()
	reply, err := messageClient.SendFriendMessage(ctx, &businesspb.SendFriendMessageRequest{
		UserId:  2,
		Content: []byte("hello gim"),
	}, grpc.Header(&header))
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("push to user", "reply", reply, "request_id", getRequestID(header))

	time.Sleep(1 * time.Second)
	fmt.Println()
	groupReply, err := messageClient.SendGroupMessage(ctx, &businesspb.SendGroupMessageRequest{
		GroupId: 1,
		Content: []byte("hello gim from group"),
	}, grpc.Header(&header))
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("push to group", "reply", groupReply, "request_id", getRequestID(header))

	time.Sleep(1 * time.Second)
	fmt.Println()
	_, err = messageClient.SendRoomMessage(ctx, &businesspb.SendRoomMessageRequest{
		RoomId:  1,
		Content: []byte("hello gim from room"),
	}, grpc.Header(&header))
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("push to room", "request_id", getRequestID(header))

	select {}
}

func generateCtx(reply *businesspb.SignInReply) context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		md.UserID:   strconv.FormatUint(reply.UserId, 10),
		md.DeviceID: strconv.FormatUint(reply.DeviceId, 10),
		md.Token:    reply.Token,
	}))
}

func getRequestID(header metadata.MD) string {
	ids := header.Get(md.RequestID)
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}
