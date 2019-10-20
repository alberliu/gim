package main

import (
	"fmt"
	"goim/test/client"
)

func main() {
	TestClient()
}

func TestClient() {
	client := client.TcpClient{}
	fmt.Println("input AppId,UserId,DeviceId,SyncSequence")
	fmt.Scanf("%d %d %d %d", &client.AppId, &client.UserId, &client.DeviceId, &client.Seq)
	client.Start()
	select {}

}
