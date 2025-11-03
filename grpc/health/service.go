package grpchealthservice

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	grpcpkg "github.com/gaiaz-iusipov/go-common/grpc"
)

var _ grpcpkg.Service = (*Service)(nil)

func New() Service {
	return Service{
		healthServer: health.NewServer(),
	}
}

type Service struct {
	healthServer *health.Server
}

func (Service) Desc() *grpc.ServiceDesc {
	return &healthpb.Health_ServiceDesc
}

func (s Service) Impl() any {
	return s.healthServer
}

func (s Service) SetServingStatus(status bool) {
	servingStatus := healthpb.HealthCheckResponse_NOT_SERVING
	if status {
		servingStatus = healthpb.HealthCheckResponse_SERVING
	}
	s.healthServer.SetServingStatus("", servingStatus)
}
