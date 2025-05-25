// Package redisinfra provides redis-related infrastructure capabilities.
package redisinfra

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tmybsv/leadgen-test-task/internal/domain/hash"
)

// HashRepository represents Redis hash repository.
type HashRepository struct {
	redisCli *redis.Client
	ttl      time.Duration
}

// NewHashRepository creates new instance of Redis hash repository by provided
// Redis client and values TTL.
func NewHashRepository(redisCli *redis.Client, ttl time.Duration) *HashRepository {
	return &HashRepository{
		redisCli: redisCli,
		ttl:      ttl,
	}
}

// Save saves provided hash to cache.
func (r *HashRepository) Save(ctx context.Context, h *hash.Hash) error {
	key := fmt.Sprintf("%s:input:%s", h.Algorithm().String(), h.Input())
	if err := r.redisCli.Set(ctx, key, h.Hashed(), r.ttl).Err(); err != nil {
		return fmt.Errorf("cache hash: %w", err)
	}

	return nil
}

// FindByInput finds hash by input string and algorithm.
func (r *HashRepository) FindByInput(ctx context.Context, input string, alg hash.Algorithm) (*hash.Hash, error) {
	key := fmt.Sprintf("%s:input:%s", alg.String(), input)
	hashed, err := r.redisCli.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("get from cache: %w", err)
	}

	h, err := hash.New(input, hashed, alg)
	if err != nil {
		return nil, fmt.Errorf("new hash: %w", err)
	}

	return h, nil
}
