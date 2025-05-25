// Package grpcapp provides gRPC server bootstrapping.
package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/tmybsv/leadgen-test-task/internal/application"
	grpcsrv "github.com/tmybsv/leadgen-test-task/internal/presentation/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// App represents wrapper to bootstrap gRPC server.
type App struct {
	port int
	srv  *grpc.Server
	log  *slog.Logger
}

// New creates new instance of application with given port, hash service and
// logger.
//
// Configures recovery and logging gRPC interceptors and registers server.
func New(port int, hashSvc *application.HashService, log *slog.Logger) *App {
	recOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("recovered from panic", slog.Any("panic", p))
			return status.Error(codes.Internal, "internal error")
		}),
	}

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recOpts...),
		logging.UnaryServerInterceptor(interceptorLogger(log), logOpts...),
	))

	grpcsrv.Register(srv, hashSvc)

	return &App{
		port: port,
		srv:  srv,
		log:  log,
	}
}

// Run runs a gRPC server on listened port. Returns an error if port already
// binded.
func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("listen TCP: %w", err)
	}

	a.log.Info("gRPC server starting", slog.Int("port", a.port))
	if err := a.srv.Serve(l); err != nil {
		return fmt.Errorf("serve gRPC server: %w", err)
	}

	return nil
}

// Stop stops a gRPC server gracefully or force stop if context was canceled.
func (a *App) Stop() {
	a.log.Info("gRPC server stopping", slog.Int("port", a.port))
	a.srv.GracefulStop()
}

func interceptorLogger(log *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
