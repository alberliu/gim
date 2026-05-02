package proxy

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"

	"gim/pkg/gen/proto/businesspb"
	"gim/pkg/gen/proto/logicpb"
	"gim/pkg/local"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func TestGrpcWebAndGRPC(t *testing.T) {
	t.Run("grpc", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := grpc.NewClient(local.BusinessServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		reply, err := businesspb.NewUserExtServiceClient(conn).SignIn(ctx, newSignInRequest("grpc"))
		if err != nil {
			t.Fatal(err)
		}
		assertSignInReply(t, reply)
	})

	t.Run("grpc-web", func(t *testing.T) {
		reply := grpcWebSignIn(t, newSignInRequest("grpc-web"))
		assertSignInReply(t, reply)
	})
}

func newSignInRequest(source string) *businesspb.SignInRequest {
	return &businesspb.SignInRequest{
		PhoneNumber: fmt.Sprintf("%d", time.Now().UnixNano()),
		Code:        "0",
		Device: &logicpb.Device{
			Type:          logicpb.DeviceType_DT_WEB,
			Brand:         source,
			Model:         "proxy-client-test",
			SystemVersion: "test",
			SdkVersion:    "test",
		},
	}
}

func grpcWebSignIn(t *testing.T, request *businesspb.SignInRequest) *businesspb.SignInReply {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := connect.NewClient[businesspb.SignInRequest, businesspb.SignInReply](
		&http.Client{Timeout: 5 * time.Second},
		"http://"+local.BusinessServerAddr+businesspb.UserExtService_SignIn_FullMethodName,
		connect.WithGRPCWeb(),
	)
	reply, err := client.CallUnary(ctx, connect.NewRequest(request))
	if err != nil {
		t.Fatal(err)
	}
	return reply.Msg
}

func assertSignInReply(t *testing.T, reply *businesspb.SignInReply) {
	t.Helper()

	if reply.GetUserId() == 0 {
		t.Fatal("user_id is empty")
	}
	if reply.GetDeviceId() == 0 {
		t.Fatal("device_id is empty")
	}
	if reply.GetToken() == "" {
		t.Fatal("token is empty")
	}
}
