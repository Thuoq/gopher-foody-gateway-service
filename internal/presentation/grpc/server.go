package grpc

import (
	"gopher-gateway-service/internal/config"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(cfg *config.Config, logger *zap.Logger) *grpc.Server {
	// Add interceptors here if needed (e.g. logging, auth, recovery)
	opts := []grpc.ServerOption{}

	server := grpc.NewServer(opts...)

	// TODO: Register gRPC services here

	if cfg.App.Env != "production" {
		reflection.Register(server)
	}

	return server
}
