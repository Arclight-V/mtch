package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Arclight-V/mtch/pkg/server/http/middleware"
)

const requestIDKey = "x-request-id"

func NewUnaryClientRequestIDInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		reqID, ok := middleware.RequestIDFromContext(ctx)
		if ok {
			ctx = metadata.AppendToOutgoingContext(ctx, requestIDKey, reqID)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func NewUnaryServerRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if vals := metadata.ValueFromIncomingContext(ctx, requestIDKey); len(vals) == 1 {
			ctx = middleware.NewContextWithRequestID(ctx, vals[0])
		}
		return handler(ctx, req)
	}
}
