package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	h, err := handler(ctx, req)

	status := "ok"
	if err != nil {
		status = "error"
	}

	HistogramResponseTime.WithLabelValues(status, info.FullMethod).Observe(time.Since(start).Seconds())

	return h, err
}
