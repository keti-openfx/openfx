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

func Invoke(service string, functionNamespace string, fxWatcherPort int, input []byte, timeout time.Duration) (string, error) {

	address := fmt.Sprintf("%s.%s:%d", service, functionNamespace, fxWatcherPort)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %s - %s\n", address, err.Error())
		return "", status.Error(codes.Internal, err.Error())
	}
	defer conn.Close()

	c := pb.NewFxWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, err := c.Call(ctx, &pb.Request{Input: input, Info: &pb.Info{FunctionName: service, Trigger: &pb.Trigger{Name: "grpc", Time: time.Now().UTC().String()}}})
	if err != nil {
		log.Printf("could not invoke: %s - %s\n", service, err.Error())
		if grpc.Code(err) == codes.DeadlineExceeded {
			return "", status.Error(codes.DeadlineExceeded, fmt.Sprintf("Function timed out after %fs", timeout.Seconds()))
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	return r.Output, nil
}
