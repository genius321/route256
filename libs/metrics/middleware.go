package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MiddlewareServerGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	h, err := handler(ctx, req)

	code := status.Code(err).String()

	status := "ok"
	if err != nil {
		status = "error"
	}

	CounterRequests.Add(1)

	CounterRequestsByGroup.WithLabelValues(info.FullMethod).Add(1)

	HistogramResponseServerTime.WithLabelValues(code, status, info.FullMethod).Observe(time.Since(start).Seconds())

	return h, err
}

func MiddlewareClientGRPC(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)
	code := status.Code(err).String()

	status := "ok"
	if err != nil {
		status = "error"
	}

	HistogramResponseClientTime.WithLabelValues(code, status, method).Observe(time.Since(start).Seconds())

	return err
}
