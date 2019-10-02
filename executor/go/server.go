package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"plugin"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/keti-openfx/openfx/executor/go/pb"
)

type FxWatcher struct {
	userFunction func(pb.Request) string
}

func NewFxWatcher() *FxWatcher {
	return &FxWatcher{}
}

func (s *FxWatcher) Call(ctx context.Context, in *pb.Request) (*pb.Reply, error) {

	if s.userFunction == nil {
		statusError := status.Error(codes.Internal, "before function call, fetch first.")
		return nil, statusError
	}

	output := s.userFunction(*in)
	return &pb.Reply{Output: output}, nil
}

func loadUserFunction(file, function string) func(pb.Request) string {
	p, err := plugin.Open(file)
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup(function)
	if err != nil {
		panic("Function not found")
	}
	return f.(func(pb.Request) string)
}

func main() {
	port := getEnvInt("PORT", 50051)
	handlerName := getEnvString("HANDLER_NAME", "Handler")
	handlerFilePath := getEnvString("HANDLER_FILE", "/go/src/github.com/keti-openfx/openfx/executor/go")

	fw := NewFxWatcher()
	fw.userFunction = loadUserFunction(handlerFilePath, handlerName)

	s := grpc.NewServer()
	pb.RegisterFxWatcherServer(s, fw)

	reflection.Register(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panicf("[fxwatcher] failed to listen: %v\n", err)
	}

	path, err := createLockFile()
	if err != nil {
		log.Panicf("Cannot write %s. Error: %s.\n", path, err.Error())
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM)
		<-sig
		s.GracefulStop()
		log.Printf("[fxwatcher] received SIGTERM.")
	}()

	log.Println("[fxwatcher] start service.")
	if err := s.Serve(lis); err != nil {
		log.Printf("[fxwatcher] failed to serve: %v\n", err)
	}

}

func createLockFile() (string, error) {
	path := filepath.Join(os.TempDir(), ".lock")
	log.Printf("Writing lock-file to: %s\n", path)
	err := ioutil.WriteFile(path, []byte{}, 0660)

	return path, err

}

func getEnvTime(key string, defaultValue time.Duration) time.Duration {
	v := os.Getenv(key)
	if v != "" {
		parsedVal, parseErr := strconv.Atoi(v)
		if parseErr == nil && parsedVal >= 0 {
			return time.Duration(parsedVal) * time.Second
		}
	}

	duration, durationErr := time.ParseDuration(v)
	if durationErr != nil {
		return defaultValue
	}
	return duration
}

func getEnvInt(key string, defaultValue int) int {
	res := defaultValue
	if v := os.Getenv(key); v != "" {
		intVal, err := strconv.Atoi(v)
		if err == nil {
			res = intVal
		}
	}
	return res
}

func getEnvString(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
