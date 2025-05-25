// Package app provides a bootstrapping capabilities.
package app

import (
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	grpcapp "github.com/tmybsv/leadgen-test-task/internal/app/grpc"
	"github.com/tmybsv/leadgen-test-task/internal/application"
	"github.com/tmybsv/leadgen-test-task/internal/domain/hash"
	redisinfra "github.com/tmybsv/leadgen-test-task/internal/infrastructure/cache/redis"
	"github.com/tmybsv/leadgen-test-task/internal/infrastructure/config"
	"github.com/tmybsv/leadgen-test-task/internal/infrastructure/hasher"
)

// App represents main application with gRPC server and Redis client.
type App struct {
	GRPCServer *grpcapp.App
	redisCli   *redis.Client
	log        *slog.Logger
}

// New creates new app instance with given configuration and logger.
//
// Initializes Redis client, hashes repository, hash service with MD5 and SHA256
// algorithms support and then creates gRPC server.
func New(cfg *config.Config, log *slog.Logger) *App {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		Username: cfg.Redis.Username,
	})

	hashRepo := redisinfra.NewHashRepository(redisCli, cfg.Redis.TTL)

	hashers := map[hash.Algorithm]hash.Hasher{
		hash.AlgorithmMD5:    &hasher.MD5{},
		hash.AlgorithmSHA256: &hasher.SHA256{},
	}

	hashSvc := application.NewHashService(hashRepo, hashers)

	grpcApp := grpcapp.New(cfg.GRPC.Port, hashSvc, log)

	return &App{
		GRPCServer: grpcApp,
		redisCli:   redisCli,
		log:        log,
	}
}

// Stop stops a gRPC server gracefully and closes connection with Redis.
func (a *App) Stop() error {
	a.GRPCServer.Stop()
	if err := a.redisCli.Close(); err != nil {
		return fmt.Errorf("close redis connecion: %w", err)
	}

	return nil
}
