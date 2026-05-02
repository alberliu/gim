package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	"gim/internal/proxy"
	"gim/pkg/server"
)

func main() {
	grpcServer := grpc.NewServer(
		grpc.UnknownServiceHandler(proxy.ServiceHandler),
		grpc.ForceServerCodec(proxy.BytesCodec{}),
	)

	grpcWebServer := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
		grpcweb.WithAllowedRequestHeaders([]string{"*"}),
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
	)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// grpc-web sets ProtoMajor=2 but doesn't strip HTTP/1.1 hop-by-hop headers;
		// gRPC-go v1.80 enforces A41: Connection header in HTTP/2 causes RST_STREAM PROTOCOL_ERROR.
		r.Header.Del("connection")
		r.Header.Del("keep-alive")
		r.Header.Del("proxy-connection")

		// gRPC-Web 请求
		if grpcWebServer.IsGrpcWebRequest(r) ||
			grpcWebServer.IsAcceptableGrpcCorsRequest(r) {
			grpcWebServer.ServeHTTP(w, r)
			return
		}

		// 原生 gRPC 请求：HTTP/2 + content-type: application/grpc
		if r.ProtoMajor == 2 &&
			strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
			return
		}

		http.NotFound(w, r)
	})

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	server.WaitForShutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := httpServer.Shutdown(ctx)
	if err != nil {
		slog.Error("http server shutdown failed", "error", err)
	}

	grpcServer.GracefulStop()
}
