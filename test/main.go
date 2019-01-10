package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"goim/test/client"
	"net/url"
	"time"
)

func main() {
	TestClient()
	//createWS()
}

func TestClient() {
	client := client.TcpClient{}
	fmt.Println("input UserId,DeviceId,Token,SendSequence,SyncSequence")
	fmt.Scanf("%d %d %s %d %d", &client.UserId, &client.DeviceId, &client.Token, &client.SendSequence, &client.SyncSequence)
	//client.UserId = 2
	//client.DeviceId = 2
	//client.Token = "1e52ec1f-3ef5-4d12-a9c6-c3b3eb995fae"
	//client.SendSequence = 1
	//client.SyncSequence = 1
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

var addr = flag.String("addr", "0.0.0.0:8888", "http service address")

func createWS() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 2)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}
