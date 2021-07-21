package main

import (
	"context"
	"fmt"
	sdk "github.com/keti-openfx/openfx/executor/go/pb"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	address := fmt.Sprintf("10.233.70.188:50051")

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %s\n", "sorry")
	}
	defer conn.Close()

	c := sdk.NewFxWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Call(ctx, &sdk.Request{Input: []byte("Invoke through API Gateway and Virtual machine"), Info: &sdk.Info{Trigger: &sdk.Trigger{Name: "grpc", Time: time.Now().UTC().String()}}})
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	log.Printf(r.Output)
}
