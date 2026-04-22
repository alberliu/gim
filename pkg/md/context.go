package md

import (
	"context"
	"net"
	"strconv"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	UserID    = "user_id"
	DeviceID  = "device_id"
	Token     = "token"
	RequestID = "request_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(RequestID, requestID))
}

func Get(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values, ok := md[key]
	if !ok || len(values) == 0 {
		return ""
	}
	return values[0]
}

// GetRequestID 获取ctx的app_id
func GetRequestID(ctx context.Context) string {
	return Get(ctx, RequestID)
}

func GetUserID(ctx context.Context) uint64 {
	userID, _ := strconv.ParseUint(Get(ctx, UserID), 10, 64)
	return userID
}

func GetDeviceID(ctx context.Context) uint64 {
	deviceID, _ := strconv.ParseUint(Get(ctx, DeviceID), 10, 64)
	return deviceID
}

// GetToken 获取ctx的token
func GetToken(ctx context.Context) string {
	return Get(ctx, Token)
}

// NewAndCopyRequestID 创建一个context,并且复制RequestID
func NewAndCopyRequestID(ctx context.Context) context.Context {
	requestID := GetRequestID(ctx)
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(RequestID, requestID))
}

func GetClientIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := p.Addr.(*net.TCPAddr); ok {
			return tcpAddr.IP.String()
		}
	}
	return ""
}
