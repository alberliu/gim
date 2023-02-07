package room

import (
	"fmt"
	"gim/pkg/protocol/pb"
	"gim/pkg/util"
	"testing"
	"time"
)

func Test_roomService_DelExpireMessage(t *testing.T) {
	err := RoomService.DelExpireMessage(1)
	fmt.Println(err)
}

func Test_roomService_List(t *testing.T) {
	msgs, err := RoomMessageRepo.List(1, 1)
	fmt.Println(err)
	fmt.Println(msgs)
}

func Test_roomService_AddMessage(t *testing.T) {
	for i := 1; i <= 20; i++ {
		err := RoomService.AddMessage(1, &pb.Message{
			Seq:      int64(i),
			SendTime: util.UnixMilliTime(time.Now()),
		})
		fmt.Println(i, err)
		time.Sleep(time.Second)
	}

	err := RoomService.DelExpireMessage(1)
	fmt.Println(err)
}
