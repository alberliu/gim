package test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gim/pkg/logger"
	"gim/pkg/md"
	"gim/pkg/protocol/pb/connectpb"
	pb "gim/pkg/protocol/pb/logicpb"
)

func TestClient(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: logger.ReplaceAttr,
	})))
	var header metadata.MD

	initData()

	network := NetworkWebsocket
	connect(network, 1, 11)
	connect(network, 1, 12)
	connect(network, 2, 2)
	connect(network, 3, 3)

	time.Sleep(2 * time.Second)
	fmt.Println()
	reply, err := getMessageIntClient().PushToUsers(context.Background(), &pb.PushToUsersRequest{
		UserIds:   []uint64{1},
		Command:   connectpb.MessageCommand_MC_USER_MESSAGE,
		Content:   []byte("hello gim"),
		IsPersist: true,
	}, grpc.Header(&header))
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("push to user", "message_id", reply.MessageId, "request_id", getRequestID(header))

	time.Sleep(1 * time.Second)
	fmt.Println()
	groupReply, err := getGroupIntClient().Push(context.Background(), &pb.GroupPushRequest{
		GroupId:   1,
		Command:   connectpb.MessageCommand_MC_GROUP_MESSAGE,
		Content:   []byte("hello gim from group"),
		IsPersist: true,
	}, grpc.Header(&header))
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("push to group", "message_id", groupReply.MessageId, "request_id", getRequestID(header))

	time.Sleep(1 * time.Second)
	fmt.Println()
	_, err = getRoomIntClient().PushRoom(context.Background(), &pb.PushRoomRequest{
		RoomId:     1,
		Command:    10000,
		Content:    []byte("hello gim from room"),
		SendTime:   time.Now().Unix(),
		IsPriority: false,
	}, grpc.Header(&header))
	if err != nil {
		t.Fatal(err)
	}
	slog.Info("push to room", "request_id", getRequestID(header))

	select {}
}

func getRequestID(header metadata.MD) string {
	ids := header.Get(md.RequestID)
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}
