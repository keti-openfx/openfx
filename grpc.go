package main

import (
	"context"

	"github.com/keti-openfx/openfx/pb"
	"google.golang.org/grpc"
)

// -----------------------------------------------------------------------------

func prepareGRPC(context context.Context, server *FxGateway) (*grpc.Server, error) {

	grpcServer := grpc.NewServer(grpc.MaxConcurrentStreams(uint32(100)))
	pb.RegisterFxGatewayServer(grpcServer, server)

	return grpcServer, nil
}
