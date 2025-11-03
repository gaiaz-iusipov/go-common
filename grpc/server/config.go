package grpcserver

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	grpcpkg "github.com/gaiaz-iusipov/go-common/grpc"
)

type config struct {
	grpcOptions      []grpc.ServerOption
	services         []grpcpkg.Service
	enableReflection bool
}

var defaultConfig = config{
	grpcOptions: []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	},
	enableReflection: true,
}

type Option func(cfg *config)

func WithGRPCOptions(grpcOptions ...grpc.ServerOption) Option {
	return func(cfg *config) {
		cfg.grpcOptions = append(cfg.grpcOptions, grpcOptions...)
	}
}

func WithServices(services ...grpcpkg.Service) Option {
	return func(cfg *config) {
		cfg.services = append(cfg.services, services...)
	}
}

func WithReflection(enableReflection bool) Option {
	return func(cfg *config) {
		cfg.enableReflection = enableReflection
	}
}
