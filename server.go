package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/keti-openfx/openfx-gateway/config"
	"github.com/keti-openfx/openfx-gateway/metrics"
	"github.com/keti-openfx/openfx-gateway/pb"
	"github.com/keti-openfx/openfx-gateway/service"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

type FxGateway struct {
	conf           config.FxGatewayConfig
	httpServer     *http.Server
	grpcServer     *grpc.Server
	kubeClient     *kubernetes.Clientset
	metricsOptions metrics.MetricOptions
	metricsFetcher metrics.PrometheusQueryFetcher
}

func NewFxGateway(c config.FxGatewayConfig, k *kubernetes.Clientset) *FxGateway {
	return &FxGateway{
		conf:           c,
		kubeClient:     k,
		metricsOptions: metrics.BuildMetricsOptions(),
		metricsFetcher: metrics.NewPrometheusQuery(c.PrometheusHost, c.PrometheusPort, &http.Client{}),
	}
}

func (f *FxGateway) Invoke(c context.Context, s *pb.InvokeServiceRequest) (*pb.Message, error) {
	start := time.Now()
	output, err := service.Invoke(s.Service, f.conf.FunctionNamespace, f.conf.FxWatcherPort, s.Input, f.conf.InvokeTimeout)
	end := time.Since(start)
	if err != nil {
		return nil, err
	}
	// For Monitoring /////////////////////////////////////////////////////////
	// function이 호출될 때마다 invoke count, time 정보 수집
	f.metricsOptions.Notify(s.Service, end, "OK")
	//////////////////////////////////////////////////////////////////////////
	return &pb.Message{Msg: output}, nil
}
func (f *FxGateway) List(c context.Context, s *pb.Empty) (*pb.Functions, error) {
	functions, err := service.List(f.conf.FunctionNamespace, f.kubeClient)
	if err != nil {
		return nil, err
	}
	// For Monitoring /////////////////////////////////////////////////////////
	// function 정보 요청이 들어오면, prometheus 서버에 쿼리를 보내서 수집한 매트릭 정보를 가져옴
	fns := metrics.AddMetricsFunctions(functions, f.metricsFetcher)
	//////////////////////////////////////////////////////////////////////////
	return &pb.Functions{Functions: fns}, nil
}
func (f *FxGateway) Deploy(c context.Context, s *pb.CreateFunctionRequest) (*pb.Message, error) {
	deployConfig := &service.DeployHandlerConfig{
		EnableHttpProbe:   f.conf.EnableHttpProbe,
		ImagePullPolicy:   f.conf.ImagePullPolicy,
		FunctionNamespace: f.conf.FunctionNamespace,
		FxWatcherPort:     f.conf.FxWatcherPort,
		SecretMountPath:   f.conf.SecretMountPath,
	}
	err := service.Deploy(s, f.kubeClient, deployConfig)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}
func (f *FxGateway) Delete(c context.Context, s *pb.DeleteFunctionRequest) (*pb.Message, error) {
	err := service.Delete(s.FunctionName, f.conf.FunctionNamespace, f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}
func (f *FxGateway) Update(c context.Context, s *pb.CreateFunctionRequest) (*pb.Message, error) {
	err := service.Update(f.conf.FunctionNamespace, s, f.kubeClient, f.conf.SecretMountPath)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}
func (f *FxGateway) GetMeta(c context.Context, s *pb.FunctionRequest) (*pb.Function, error) {
	fn, err := service.GetMeta(s.FunctionName, f.conf.FunctionNamespace, f.kubeClient)
	if err != nil {
		return nil, err
	}
	// For Monitoring /////////////////////////////////////////////////////////
	// function 정보 요청이 들어오면, prometheus 서버에 쿼리를 보내서 수집한 매트릭 정보를 가져옴
	fn = metrics.AddMetricsFunction(fn, f.metricsFetcher)
	///////////////////////////////////////////////////////////////////////////
	return fn, nil
}
func (f *FxGateway) GetLog(c context.Context, s *pb.FunctionRequest) (*pb.Message, error) {
	log, err := service.GetLog(s.FunctionName, f.conf.FunctionNamespace, f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: log}, nil
}
func (f *FxGateway) ReplicaUpdate(c context.Context, s *pb.ScaleServiceRequest) (*pb.Message, error) {
	err := service.ReplicaUpdate(f.conf.FunctionNamespace, s, f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}
func (f *FxGateway) Info(c context.Context, s *pb.Empty) (*pb.Message, error) {
	info, err := service.Info(f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: info}, nil
}
func (f *FxGateway) HealthCheck(c context.Context, s *pb.Empty) (*pb.Message, error) {
	return &pb.Message{Msg: "OK"}, nil
}

// -----------------------------------------------------------------------------

// Start Openfx Gateway
func (f *FxGateway) Start() error {

	// For Monitoring /////////////////////////////////////////////////////////
	//
	//
	// 매트릭 정보를 수집하는 Exporter를 생성
	exporter := metrics.NewExporter(f.metricsOptions)
	// 5초마다 function들의 Replica(복제본 수) 정보  수집
	servicePollInterval := time.Second * 5
	exporter.StartServiceWatcher(f.conf.FunctionNamespace, f.kubeClient, f.metricsOptions, servicePollInterval)
	// Prometheus 매트릭 수집기 등록
	metrics.RegisterExporter(exporter)

	///////////////////////////////////////////////////////////////////////////

	var err error

	// Initialize listener
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", f.conf.TCPPort))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// tcpMuxer
	tcpMux := cmux.New(conn)

	// Connection dispatcher rules
	grpcL := tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := tcpMux.Match(cmux.HTTP1Fast())

	// initialize gRPC server instance
	f.grpcServer, err = prepareGRPC(ctx, f)
	if err != nil {
		log.Fatalln("Unable to initialize gRPC server instance")
		return err
	}

	// initialize HTTP server
	f.httpServer, err = prepareHTTP(ctx, fmt.Sprintf("localhost:%d", f.conf.TCPPort), f.conf.FunctionNamespace, f.conf.FxWatcherPort, f.conf.InvokeTimeout, f.conf.ReadTimeout, f.conf.WriteTimeout, f.conf.IdleTimeout)
	if err != nil {
		log.Fatalln("Unable to initialize HTTP server instance")
		return err
	}

	// Start servers
	go func() {
		if err := f.grpcServer.Serve(grpcL); err != nil {
			log.Fatalln("Unable to start external gRPC server")
		}
	}()
	go func() {
		if err := f.httpServer.Serve(httpL); err != nil {
			log.Fatalln("Unable to start HTTP server")
		}
	}()

	return tcpMux.Serve()
}
