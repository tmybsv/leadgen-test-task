// Package application provides core business logic and infrastructure
// interaction.
package application

import (
	"context"
	"fmt"

	"github.com/tmybsv/leadgen-test-task/internal/domain/hash"
)

// HashService serves hash business logic. Contains implementation of hash
// repository and map of hashers.
type HashService struct {
	hashRepo hash.Repository
	hashers  map[hash.Algorithm]hash.Hasher
}

// NewHashService creates new instance of hash service.
func NewHashService(hashRepo hash.Repository, hashers map[hash.Algorithm]hash.Hasher) *HashService {
	return &HashService{
		hashRepo: hashRepo,
		hashers:  hashers,
	}
}

// CreateHash creates hash of provided string by given algorithm.
//
// Uses a cache-first approach. Only if hash string not found in cache will
// create a new one.
func (s *HashService) CreateHash(ctx context.Context, input string, alg hash.Algorithm) (*hash.Hash, error) {
	h, err := s.hashRepo.FindByInput(ctx, input, alg)
	if err == nil {
		return h, nil
	}

	hasher, ok := s.hashers[alg]
	if !ok {
		return nil, fmt.Errorf("hasher for algorithm %v not registered", alg)
	}

	hashed := hasher.Hash(input)
	h, err = hash.New(input, hashed, alg)
	if err != nil {
		return nil, fmt.Errorf("new hash: %w", err)
	}

	if err = s.hashRepo.Save(ctx, h); err != nil {
		return nil, fmt.Errorf("save hash for %q: %w", input, err)
	}

	return h, nil
}
