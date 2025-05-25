// Package grpcsrv provides gRPC server.
package grpcsrv

import (
	"context"
	"errors"

	"github.com/tmybsv/leadgen-test-task/internal/application"
	"github.com/tmybsv/leadgen-test-task/internal/domain/hash"
	pbhasher "github.com/tmybsv/leadgen-test-task/pkg/pb/hasher/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type hashServer struct {
	pbhasher.UnimplementedHasherServiceServer
	hashSvc *application.HashService
}

// Register wraps a native gRPC register and registers gRPC server
// implementation.
func Register(s *grpc.Server, hashSvc *application.HashService) {
	pbhasher.RegisterHasherServiceServer(s, &hashServer{
		hashSvc: hashSvc,
	})
}

func (s *hashServer) Hash(ctx context.Context, req *pbhasher.HashRequest) (*pbhasher.HashResponse, error) {
	if req.Input == "" {
		return nil, status.Error(codes.InvalidArgument, "input is required")
	}

	domainAlg, err := convertAlgorithm(req.Algorithm)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	h, err := s.hashSvc.CreateHash(ctx, req.Input, domainAlg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbhasher.HashResponse{
		Hash: h.Hashed(),
	}, nil
}

func convertAlgorithm(pbAlg pbhasher.HashAlgorithm) (hash.Algorithm, error) {
	switch pbAlg {
	case pbhasher.HashAlgorithm_HASH_ALGORITHM_MD5:
		return hash.AlgorithmMD5, nil
	case pbhasher.HashAlgorithm_HASH_ALGORITHM_SHA256:
		return hash.AlgorithmSHA256, nil
	default:
		return 0, errors.New("unsupported algorithm")
	}
}
