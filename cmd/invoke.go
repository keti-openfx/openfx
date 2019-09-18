package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/keti-openfx/openfx/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Invoke(service string, functionNamespace string, fxWatcherPort int, input []byte, timeout time.Duration) (string, error) {

	address := fmt.Sprintf("%s.%s:%d", service, functionNamespace, fxWatcherPort)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return "", status.Error(codes.Internal, "did not connect: "+address+" - "+err.Error())
	}
	defer conn.Close()

	c := pb.NewFxWatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	r, err := c.Call(ctx, &pb.Request{Input: input, Info: &pb.Info{FunctionName: service, Trigger: &pb.Trigger{Name: "grpc", Time: time.Now().UTC().String()}}})
	if err != nil {
		if grpc.Code(err) == codes.DeadlineExceeded {
			return "", status.Error(codes.DeadlineExceeded, fmt.Sprintf("Function timed out after %fs", time.Since(start).Seconds()))
		}
		return "", status.Error(codes.Internal, "could not invoke: "+service+" - "+err.Error())
	}

	return r.Output, nil
}
