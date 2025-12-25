package main

import (
	cachev1 "cache-service/gen/cache/v1"
	"cache-service/internal/cache"
	"cache-service/internal/grpcserver"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultGrpcAddr = ":50051"

	gracefulStopTimeout = 10 * time.Second
)

func main() {
	log.Println("Cache service started")

	addr := getEnv("GRPC_ADDR", defaultGrpcAddr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s is failed, err %v\n", err)
	}

	srv := cache.NewCacheServiceImpl()
	grpcHandler := grpcserver.New(srv)

	grpcSrv := grpc.NewServer()

	cachev1.RegisterCacheServiceServer(grpcSrv, grpcHandler)

	reflection.Register(grpcSrv)

	errChan := make(chan error, 1)
	go func() {
		log.Printf("gRPC server listening on %s", addr)
		if err := grpcSrv.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		log.Printf("shutdown signal received: %v", ctx.Err())
	case err := <-errChan:
		log.Printf("gRPC server stopped with error: %v", err)
	}

	gracefulStop(grpcSrv, gracefulStopTimeout)

	log.Printf("server stopped")
}


func gracefulStop(s *grpc.Server, timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		log.Printf("graceful stop timeout (%s), forcing stop", timeout)
		s.Stop()
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
