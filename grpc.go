package main

import (
	"context"

	"github.com/keti-openfx/openfx-gateway/pb"
	"google.golang.org/grpc"
)

// -----------------------------------------------------------------------------

func prepareGRPC(context context.Context, server *FxGateway) (*grpc.Server, error) {

	grpcServer := grpc.NewServer()
	pb.RegisterFxGatewayServer(grpcServer, server)

	return grpcServer, nil
}
