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
	fmt.Println("input UserId,DeviceId,Token,SendSequence,SyncSequence")
	fmt.Scanf("%d %d %s %d %d", &client.UserId, &client.DeviceId, &client.Token, &client.SendSequence, &client.SyncSequence)
	client.Start()
	for {
		client.SendMessage()
	}
}

func TestDebugClient() {
	client := client.TcpClient{}
	fmt.Println("input UserId,DeviceId,Token,SendSequence,SyncSequence")
	fmt.Scanf("%d %d %s %d %d", &client.UserId, &client.DeviceId, &client.Token, &client.SendSequence, &client.SyncSequence)
	client.Start()
	for {
		client.SendMessage()
	}
}
