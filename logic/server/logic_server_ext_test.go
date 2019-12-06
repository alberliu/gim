package server

import (
	"context"
	"gim/public/pb"
	"reflect"
	"testing"
)

func TestLogicServerExtServer_SendMessage(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *pb.SendMessageReq
	}
	tests := []struct {
		name    string
		l       *LogicServerExtServer
		args    args
		want    *pb.SendMessageResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LogicServerExtServer{}
			got, err := l.SendMessage(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogicServerExtServer.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogicServerExtServer.SendMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
