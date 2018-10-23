package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/keti-openfx/openfx-gateway/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Invoke(service string, functionNamespace string, fxWatcherPort int, input []byte) (string, error) {

	address := fmt.Sprintf("%s.%s:%d", service, functionNamespace, fxWatcherPort)
	log.Printf("Invoke Service Address:%s\n", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %s\n", address)
		return "", status.Error(codes.Internal, err.Error())
	}
	defer conn.Close()

	c := pb.NewFxWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Call(ctx, &pb.Request{Input: input, Info: &pb.Info{FunctionName: service, Trigger: &pb.Trigger{Name: "grpc", Time: time.Now().UTC().String()}}})
	if err != nil {
		log.Printf("could not invoke: %s\n", service)
		return "", status.Error(codes.Internal, err.Error())
	}
	log.Printf("FxWatcher Invoke Output: %s", r.Output)

	return r.Output, nil
}
