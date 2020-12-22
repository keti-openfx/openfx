package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/keti-openfx/fx-idler/pb"
	grpcgo "google.golang.org/grpc"
)

func GetMeta(functionName, fxGateway string) (*pb.Function, error) {

	gateway := strings.TrimRight(fxGateway, "/")

	conn, err := grpcgo.Dial(gateway, grpcgo.WithInsecure())
	if err != nil {
		return nil, errors.New("did not connect: " + err.Error())
	}

	client := pb.NewFxGatewayClient(conn)

	function, statusErr := client.GetMeta(context.Background(), &pb.FunctionRequest{FunctionName: functionName})
	if statusErr != nil {
		return nil, errors.New("did not get meta: " + statusErr.Error())
	}

	return function, nil
}
