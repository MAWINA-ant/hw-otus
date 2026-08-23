package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// loggingInterceptor logs every unary RPC call, mirroring the access-log
// format used by the HTTP server's logging middleware.
func (s *Server) loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	latency := time.Since(start)
	code := status.Code(err)

	logLine := fmt.Sprintf("%s [%s] %s %d %d",
		clientIP(ctx),
		start.Format("02/Jan/2006:15:04:05 -0700"),
		info.FullMethod,
		code,
		latency.Milliseconds(),
	)

	if err != nil {
		s.logger.Info(logLine + " error=\"" + err.Error() + "\"")
	} else {
		s.logger.Info(logLine)
	}

	return resp, err
}

func clientIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok || p.Addr == nil {
		return "-"
	}
	return p.Addr.String()
}
