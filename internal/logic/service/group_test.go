package service

import (
	"context"
	"gim/internal/logic/model"
	"gim/pkg/pb"
	"reflect"
	"testing"
)

func Test_groupService_Create(t *testing.T) {
	type args struct {
		ctx       context.Context
		userId    int64
		group     model.Group
		memberIds []int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &groupService{}
			got, err := gr.Create(tt.args.ctx, tt.args.userId, tt.args.group, tt.args.memberIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupService_Get(t *testing.T) {
	type args struct {
		ctx     context.Context
		groupId int64
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Group
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &groupService{}
			got, err := gr.Get(tt.args.ctx, tt.args.groupId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupService_GetUsers(t *testing.T) {
	type args struct {
		ctx     context.Context
		groupId int64
	}
	tests := []struct {
		name    string
		args    args
		want    []*pb.GroupMember
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &groupService{}
			got, err := s.GetUsers(tt.args.ctx, tt.args.groupId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupService_Update(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId int64
		group  model.Group
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &groupService{}
			if err := gr.Update(tt.args.ctx, tt.args.userId, tt.args.group); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
