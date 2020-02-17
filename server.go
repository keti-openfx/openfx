package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/keti-openfx/openfx/cmd"
	"github.com/keti-openfx/openfx/config"
	"github.com/keti-openfx/openfx/metrics"
	"github.com/keti-openfx/openfx/pb"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

/* FxGateway는 gRPC service 구현체 */
type FxGateway struct {
	conf       config.FxGatewayConfig /* 환경변수를 통한 설정 */
	kubeClient *kubernetes.Clientset  /* kubernetes 클라이언트 */

	httpServer *http.Server
	grpcServer *grpc.Server

	metricsOptions metrics.MetricOptions          /* prometheus 메트릭 */
	metricsFetcher metrics.PrometheusQueryFetcher /* prometheus 클라이언트 */

	events trace.EventLog
}

// FxGateway 생성
func NewFxGateway(c config.FxGatewayConfig, k *kubernetes.Clientset) *FxGateway {
	gw := &FxGateway{
		conf:           c,
		kubeClient:     k,
		metricsOptions: metrics.BuildMetricsOptions(),
		metricsFetcher: metrics.NewPrometheusQuery(c.PrometheusHost, c.PrometheusPort, &http.Client{}),
	}

	if EnableTracing {
		_, file, line, _ := runtime.Caller(1)
		gw.events = trace.NewEventLog("FxGateway", fmt.Sprintf("%s:%d", file, line))
	}

	return gw
}

/* grpc handler
 * function 호출*/
func (f *FxGateway) Invoke(c context.Context, s *pb.InvokeServiceRequest) (*pb.Message, error) {
	start := time.Now()
	output, err := cmd.Invoke(s.Service, f.conf.FunctionNamespace, f.conf.FxWatcherPort, s.Input, f.conf.InvokeTimeout)
	end := time.Since(start)
	if err != nil {
		return nil, err
	}
	// For Monitoring /////////////////////////////////////////////////////////
	// function이 호출될 때마다 invoke count를 증가와 invoke에 걸린 시간 정보 수집
	f.metricsOptions.Notify(s.Service, end, "OK")
	//////////////////////////////////////////////////////////////////////////
	return &pb.Message{Msg: output}, nil
}

// grpc handler
// function list 조회
func (f *FxGateway) List(c context.Context, s *pb.Empty) (*pb.Functions, error) {
	functions, err := cmd.List(f.conf.FunctionNamespace, f.kubeClient)
	if err != nil {
		return nil, err
	}
	// For Monitoring /////////////////////////////////////////////////////////
	// function 정보 요청이 들어오면, prometheus 서버에 쿼리를 보내서 수집한 매트릭 정보를 가져옴
	fns := metrics.AddMetricsFunctions(functions, f.metricsFetcher)
	//////////////////////////////////////////////////////////////////////////
	return &pb.Functions{Functions: fns}, nil
}

// grpc handler
// function 배포
func (f *FxGateway) Deploy(c context.Context, s *pb.CreateFunctionRequest) (*pb.Message, error) {
	deployConfig := &cmd.DeployHandlerConfig{
		EnableHttpProbe:   f.conf.EnableHttpProbe,
		ImagePullPolicy:   f.conf.ImagePullPolicy,
		FunctionNamespace: f.conf.FunctionNamespace,
		FxWatcherPort:     f.conf.FxWatcherPort,
		SecretMountPath:   f.conf.SecretMountPath,
	}
	err := cmd.Deploy(s, f.kubeClient, deployConfig)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}

// grpc handler
// function 삭제
func (f *FxGateway) Delete(c context.Context, s *pb.DeleteFunctionRequest) (*pb.Message, error) {
	err := cmd.Delete(s.FunctionName, f.conf.FunctionNamespace, f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}

// grpc handler
// function 업데이트
func (f *FxGateway) Update(c context.Context, s *pb.CreateFunctionRequest) (*pb.Message, error) {
	err := cmd.Update(f.conf.FunctionNamespace, s, f.kubeClient, f.conf.SecretMountPath)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}

// grpc handler
// function 정보 조회
func (f *FxGateway) GetMeta(c context.Context, s *pb.FunctionRequest) (*pb.Function, error) {
	fn, err := cmd.GetMeta(s.FunctionName, f.conf.FunctionNamespace, f.kubeClient)
	if err != nil {
		return nil, err
	}
	// For Monitoring /////////////////////////////////////////////////////////
	// function 정보 요청이 들어오면, prometheus 서버에 쿼리를 보내서 수집한 매트릭 정보를 가져옴
	fn = metrics.AddMetricsFunction(fn, f.metricsFetcher)
	///////////////////////////////////////////////////////////////////////////
	return fn, nil
}

// grpc handler
// function의 출력과 에러 조회
func (f *FxGateway) GetLog(c context.Context, s *pb.FunctionRequest) (*pb.Message, error) {
	log, err := cmd.GetLog(s.FunctionName, f.conf.FunctionNamespace, f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: log}, nil
}

// grpc handler
// function의 복제본 수 업데이트
func (f *FxGateway) ReplicaUpdate(c context.Context, s *pb.ScaleServiceRequest) (*pb.Message, error) {
	err := cmd.ReplicaUpdate(f.conf.FunctionNamespace, s, f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}

// grpc handler
// gateway의 버전 정보 조회
func (f *FxGateway) Info(c context.Context, s *pb.Empty) (*pb.Message, error) {
	info, err := cmd.Info(f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: info}, nil
}

// grpc handler
// gateway의 health check
func (f *FxGateway) HealthCheck(c context.Context, s *pb.Empty) (*pb.Message, error) {
	return &pb.Message{Msg: "OK"}, nil
}

// -----------------------------------------------------------------------------

/*
 * Start Openfx Gateway
 * 멀티플렉서 생성, 핸들러 등록, grpc/http 서버 시작
 */
func (f *FxGateway) Start() error {

	/* For Monitoring **********************************************************/
	/*
		/* 매트릭 정보를 수집하는 Exporter를 생성 */
	exporter := metrics.NewExporter(f.metricsOptions)
	// 5초마다 function들의 Replica(복제본 수) 정보  수집
	servicePollInterval := time.Second * 5
	exporter.StartServiceWatcher(f.conf.FunctionNamespace, f.kubeClient, f.metricsOptions, servicePollInterval)
	// Prometheus 매트릭 수집기 등록
	metrics.RegisterExporter(exporter)
	///////////////////////////////////////////////////////////////////////////

	var err error

	// Initialize listener
	// 프로토콜, IP 주소, 포트 번호를 설정하여 네트워크 연결 대기
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", f.conf.TCPPort))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Multiplexer 생성
	// payload에 따라 연결 다중화, 동일한 TCP listener에서 다양한 프로토콜을 사용 가능
	tcpMux := cmux.New(conn)

	// Connection dispatcher rules
	grpcL := tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := tcpMux.Match(cmux.HTTP1Fast())

	// http/grpc server의 값, 시그널, cancelation, deadline 등을 전달하기 위해서 사용
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// initialize gRPC server instance
	// gRPC service 구현체인 FxGateway를 전달하여 gRPC server에 handler를 등록하고
	// gRPC server 반환
	f.grpcServer, err = prepareGRPC(ctx, f)
	if err != nil {
		log.Fatalln("Unable to initialize gRPC server instance")
		return err
	}

	// initialize HTTP server
	// grpc/http handler를 http server에 등록하고, http server 반환
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

	// Start Multiplexer
	return tcpMux.Serve()
}
