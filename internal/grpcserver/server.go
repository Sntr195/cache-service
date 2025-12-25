package grpcserver

import (
	"context"
	"errors"
	"time"

	cachev1 "cache-service/gen/cache/v1"
	"cache-service/internal/cache"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Server struct {
	cachev1.UnimplementedCacheServiceServer
	srv cache.CacheService
}

func New(srv cache.CacheService) *Server {
	return &Server{srv: srv}
}

func (s *Server) Get(ctx context.Context, req *cachev1.GetRequest) (*cachev1.GetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nill")
	}
	key := req.GetKey()
	if key == "" {
		return nil, status.Error(codes.InvalidArgument, "key is empty")
	}

	val, found, err := s.srv.Get(ctx, key)
	if err != nil {
		return nil, mapErr(err)
	}

	resp := &cachev1.GetResponse{
		Found: found,
	}

	if !found {
		return resp, nil
	}

	resp.Value = cloneBytes(val)

	_ = req.GetIncludeMeta()

	return resp, nil
}

func (s *Server) Set(ctx context.Context, req *cachev1.SetRequest) (*cachev1.SetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	key := req.GetKey()
	if key == "" {
		return nil, status.Error(codes.InvalidArgument, "key is empty")
	}

	ttl, err := ttlFromProto(req.Ttl)
	if err != nil {
		return nil, err
	}

	val := cloneBytes(req.GetValue())

	evicted, serr := s.srv.Set(ctx, key, val, ttl, req.GetKeepTtl())
	if serr != nil {
		return nil, mapErr(serr)
	}

	return &cachev1.SetResponse{Evicted: evicted}, nil
}

func (s *Server) Delete(ctx context.Context, req *cachev1.DeleteRequest) (*cachev1.DeleteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	key := req.GetKey()
	if key == "" {
		return nil, status.Error(codes.InvalidArgument, "key is empty")
	}

	deleted, err := s.srv.Delete(ctx, key)
	if err != nil {
		return nil, mapErr(err)
	}
	return &cachev1.DeleteResponse{Deleted: deleted}, nil
}

func (s *Server) Len(ctx context.Context, _ *cachev1.LenRequest) (*cachev1.LenResponse, error) {
	n, err := s.srv.Len(ctx)
	if err != nil {
		return nil, mapErr(err)
	}
	return &cachev1.LenResponse{Items: n}, nil
}

func ttlFromProto(d *durationpb.Duration) (time.Duration, error) {
	if d == nil {
		return 0, nil
	}

	ttl := d.AsDuration()

	if ttl < 0 {
		return 0, status.Error(codes.InvalidArgument, "ttl must be >= 0")
	}

	return ttl, nil
}

func mapErr(err error) error {
	switch {
	case errors.Is(err, cache.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, cache.ErrTooLarge):
		return status.Error(codes.ResourceExhausted, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}

func cloneBytes(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
}
