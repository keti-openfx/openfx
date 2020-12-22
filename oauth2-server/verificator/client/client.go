package main

import (
	"context"
	"fmt"
	"log"

	"github.com/keti-openfx/openfx/verificator/pb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial(":30011", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewFxauthClient(conn)

	doUnary(c)
}

func doUnary(c pb.FxauthClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &pb.AuthRequest{
		Token: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyMjIyMjIiLCJleHAiOjE1ODkxNjk1MzQsInN1YiI6IjEyMyJ9.nMfBpgDt4gNhvwj-gr3Xhu9mEHlyRKFsOScfAFAFRuVdlfJ-8p6nir-JISy0ICYjDLJfOhKNBk_fc8v7YklXig",
	}

	res, err := c.Authorization(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.IsAuth)

	return
}
