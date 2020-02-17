package client

import (
	"context"
	"log"
	"time"

	"github.com/keti-openfx/openfx/pb"
	"google.golang.org/grpc"
)

func Deploy(address string, req *pb.CreateFunctionRequest, timeout time.Duration) (string, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewFxGatewayClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, statusErr := client.Deploy(ctx, req)
	return r.Message, statusErr
}
