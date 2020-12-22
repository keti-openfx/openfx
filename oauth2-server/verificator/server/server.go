package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/keti-openfx/openfx/verificator/pb"
	"google.golang.org/grpc"
)

// Server ...
type Server struct{}

// Auth ...
func (*Server) Authentication(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	token := req.GetToken()
	fmt.Printf("check CODE %v\n", token)
	/*
		DB 체크
	*/
	res := &pb.AuthResponse{
		IsAuth: true,
	}
	return res, nil
}

// Auth ...
func (*Server) Authorization(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	token := req.GetToken()
	fmt.Printf("check CODE %v\n", token)
	/*
		DB 체크
	*/
	res := &pb.AuthResponse{
		IsAuth: true,
	}
	return res, nil
}

func main() {
	fmt.Println("Verificator Server Start..")

	lis, err := net.Listen("tcp", ":30011")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFxauthServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
