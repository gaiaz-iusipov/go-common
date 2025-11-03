package grpc

import "google.golang.org/grpc"

type Server interface {
	grpc.ServiceRegistrar
}

type Service interface {
	Desc() *grpc.ServiceDesc
	Impl() any
}
