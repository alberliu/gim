package device

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gim/pkg/protocol/pb/logicpb"
)

func getClient() pb.DeviceExtServiceClient {
	conn, err := grpc.NewClient("127.0.0.1:8010", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewDeviceExtServiceClient(conn)
}

func TestDeviceExtService_RegisterDevice(t *testing.T) {
	reply, err := getClient().RegisterDevice(context.TODO(), &pb.RegisterDeviceRequest{
		Type:          pb.DeviceType_DT_ANDROID,
		Brand:         "huawei",
		Model:         "huawei15",
		SystemVersion: "1.0.0",
		SdkVersion:    "1.0.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply)
}
