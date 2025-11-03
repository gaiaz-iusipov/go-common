package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcpkg "github.com/gaiaz-iusipov/go-common/grpc"
)

var _ grpcpkg.Server = (*Server)(nil)

func New(addr string, opts ...Option) Server {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	grpcServer := grpc.NewServer(cfg.grpcOptions...)

	for _, service := range cfg.services {
		grpcServer.RegisterService(service.Desc(), service.Impl())
	}

	if cfg.enableReflection {
		reflection.Register(grpcServer)
	}

	return Server{
		addr:       addr,
		grpcServer: grpcServer,
	}
}

type Server struct {
	addr       string
	grpcServer *grpc.Server
}

func (s Server) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.grpcServer.RegisterService(desc, impl)
}

func (s Server) Run(ctx context.Context) error {
	listener, err := new(net.ListenConfig).Listen(ctx, "tcp", s.addr)
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	go func() {
		if serveErr := s.grpcServer.Serve(listener); serveErr != nil {
			slog.ErrorContext(ctx, "failed to serve grpc server", "error", serveErr)
		}
	}()

	return nil
}

func (s Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}
