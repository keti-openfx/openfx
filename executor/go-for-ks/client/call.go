package client

import (
	"context"
	"log"
	"time"

	"github.com/keti-openfx/openfx/executor/go/pb"
	"google.golang.org/grpc"
)

func Call(address string, input []byte, timeout time.Duration) string {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFxWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, err := client.Call(ctx, &pb.Request{Input: input})
	if err != nil {
		log.Fatalf("could not invoke: %v\n", err)
	}
	return r.Output
}
