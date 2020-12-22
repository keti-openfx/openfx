package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/keti-openfx/fx-idler/pb"
	grpcgo "google.golang.org/grpc"
)

func List(fxGateway string) (*pb.Functions, error) {

	gateway := strings.TrimRight(fxGateway, "/")

	conn, err := grpcgo.Dial(gateway, grpcgo.WithInsecure())
	if err != nil {
		return nil, errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	functions, statusErr := client.List(context.Background(), &pb.Empty{})
	if statusErr != nil {
		return nil, errors.New("did not listed: " + statusErr.Error())
	}

	return functions, nil
}
