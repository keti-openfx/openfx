package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/keti-openfx/fx-idler/pb"
	grpcgo "google.golang.org/grpc"
)

func Scale(fxGateway, functionName, nameSpace string, replicaCount uint64) (string, error) {

	gateway := strings.TrimRight(fxGateway, "/")

	conn, err := grpcgo.Dial(gateway, grpcgo.WithInsecure())
	if err != nil {
		return "error", errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	message, statusErr := client.Scale(context.Background(), &pb.ScaleRequest{NameSpace: nameSpace, FunctionName: functionName, Replicas: replicaCount})
	if statusErr != nil {
		return "error", errors.New("did not Scale: " + statusErr.Error())
	}

	return message.Msg, nil
}

