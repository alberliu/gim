package proxy

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gim/pkg/gen/proto/businesspb"
	"gim/pkg/gerrors"
	"gim/pkg/md"
	"gim/pkg/rpc"
	"gim/pkg/ugrpc"
)

func ServiceHandler(srv any, stream grpc.ServerStream) error {
	method, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return errors.New("get method fail")
	}

	var args, reply []byte
	err := stream.RecvMsg(&args)
	if err != nil {
		return err
	}

	err = invoke(stream.Context(), method, &args, &reply)
	if err != nil {
		return err
	}

	err = stream.SendMsg(&reply)
	if err != nil {
		return err
	}
	return nil
}

var clients sync.Map

func invoke(ctx context.Context, method string, args, reply any, opts ...grpc.CallOption) error {
	fullMethod, err := ugrpc.ParseFullMethod(method)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(fullMethod.Service, "ExtService") {
		return gerrors.ErrBadRequest
	}

	client, ok := clients.Load(fullMethod.Server)
	if !ok {
		target := ugrpc.GetTarget(fullMethod.Server)
		conn, err := grpc.NewClient(target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.ForceCodec(BytesCodec{}),
			),
			ugrpc.WithDefaultServiceConfig(),
			grpc.WithChainUnaryInterceptor(
				ugrpc.ProxyLogClientInterceptor,
				authServerInterceptor,
			))
		if err != nil {
			return err
		}

		clients.Store(fullMethod.Server, conn)
		client = conn
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return client.(*grpc.ClientConn).Invoke(ctx, method, args, reply, opts...)
}

func authServerInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if _, ok := URLWhitelist[method]; !ok {
		userID := md.GetUserID(ctx)
		deviceID := md.GetDeviceID(ctx)
		token := md.GetToken(ctx)

		_, err := rpc.GetUserIntClient().Auth(ctx, &businesspb.AuthRequest{
			UserId:   userID,
			DeviceId: deviceID,
			Token:    token,
		})
		if err != nil {
			return err
		}
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

var URLWhitelist = map[string]struct{}{
	businesspb.UserExtService_SignIn_FullMethodName: {},
}
