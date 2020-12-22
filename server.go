package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
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

// oauth2 info
const (
	authServerURL = "http://10.0.0.91:30011"
)

/*
type Access_info struct {
	Client_id  string `json:"client_id"`
	Expires_in int    `json:"expires_in"`
	Scope      string `json:"scope"`
	User_id    string `json:"user_id"`
	Grade      string `json:"grade"`
}
*/

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

func IsAuth(fnList []*pb.Function, serviceName string) bool {
	for _, fn := range fnList {
		if fn.Name == serviceName {
			return true
		}
	}
	return false
}

/* grpc handler
 * function 호출*/
func (f *FxGateway) Invoke(c context.Context, s *pb.InvokeServiceRequest) (*pb.Message, error) {

	target, err := GetTokenData(s.Token)
	if err != nil {
		return nil, err
	}

	fnList, err := cmd.List(target, f.kubeClient)
	if err != nil {
		return nil, err
	}

	IsAccessble := IsAuth(fnList, s.Service)
	if !IsAccessble {
		log.Printf("[Unauthorized] Invoke / User : %v, Namespaces : %v, function : %v", target.User_id, target.Scope, s.Service)
		return &pb.Message{Msg: "This is an unauthorized function call.\n"}, nil
	}

	start := time.Now()
	// log 입력 필요
	log.Printf("[Authorized] Invoke / User : %v, Namespaces : %v, function : %v", target.User_id, target.Scope, s.Service)
	output, err := cmd.Invoke(s.Service, target.Scope, f.conf.FxWatcherPort, s.Input, f.conf.InvokeTimeout)

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

//grpc handler
//네임스페이스 생성
func (f *FxGateway) Create(c context.Context, s *pb.CreateNamespaceRequest) (*pb.Message, error) {
	err := cmd.Create(s, f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "Create Namespaces"}, nil
}

// grpc handler
// Login
func (f *FxGateway) Login(c context.Context, s *pb.LoginRequest) (*pb.LoginResponse, error) {

	target, err := GetTokenData(s.Token)
	if err != nil {
		return nil, err
	}
	log.Println("successfuly Token")

	resp, err := http.Get("http://10.0.0.91:30366/api/requestIDE-URL/" + target.Client_id + "/" + target.User_id)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("successfuly requestIDE-URL")

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	IDE_url := string(data)

	log.Println(IDE_url)

	if IDE_url != "null" { // null check
		return &pb.LoginResponse{Msg: "[Login] Succeed"}, nil
	} else { // IDE 구성 환경 구축

		deployConfig := &cmd.DeployHandlerConfig{
			EnableHttpProbe:   f.conf.EnableHttpProbe,
			ImagePullPolicy:   f.conf.ImagePullPolicy,
			FunctionNamespace: f.conf.FunctionNamespace,
			FxWatcherPort:     4000,
			SecretMountPath:   f.conf.SecretMountPath,
		}

		s.Member = target.Scope + "-" + target.User_id

		log.Println("before ")
		IDEurl, err := cmd.CreateIDE(s, f.kubeClient, deployConfig)
		if err != nil {
			return nil, err
		}

		log.Println(IDEurl)

		// URL 서버에 URL 정보 전달
		resp_url, err := http.PostForm("http://10.0.0.91:30366/api/user", url.Values{"UserID": {target.User_id}, "url": {IDEurl}, "ClientID": {target.Client_id}})
		if err != nil {
			panic(err)
		}
		defer resp_url.Body.Close()

		return &pb.LoginResponse{Msg: "[Login] Succeed"}, nil
	}
}

// grpc handler
func (f *FxGateway) StartIDE(c context.Context, s *pb.StartRequest) (*pb.StartResponse, error) {
	// 1. 토큰 권한 확인

	target, err := GetTokenData(s.Token)
	if err != nil {
		return nil, err
	}

	// 2. url 반환받기
	resp, err := http.Get("http://10.0.0.91:30366/api/requestIDE-URL/" + target.Client_id + "/" + target.User_id)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body) // byte[] "\"10.0.0.91:3030\"" ?
	if err != nil {
		panic(err)
	}

	type Ide struct {
		IDE string
	}

	Ide_info := Ide{}

	jsonErr := json.Unmarshal(data, &Ide_info)
	if jsonErr != nil {
		return &pb.StartResponse{IDE: "[Incorrect approach Error] IDE creation is required."}, nil
	}

	if Ide_info.IDE != "null" {
		return &pb.StartResponse{IDE: Ide_info.IDE}, nil
	}

	return &pb.StartResponse{IDE: "[Incorrect approach Error] IDE creation is required."}, nil
}

// grpc handler
// Login
func (f *FxGateway) ExitIDE(c context.Context, s *pb.ExitRequest) (*pb.ExitResponse, error) {
	return &pb.ExitResponse{Msg: s.Token}, nil
}

// grpc handler
// function list 조회
func (f *FxGateway) List(c context.Context, s *pb.TokenRequest) (*pb.Functions, error) {

	// //정보값 전송 후 받아오기
	target, err := GetTokenData(s.Token)
	if err != nil {
		return nil, err
	}

	if target.Scope == "" {
		return nil, errors.New("Input Valid Access Token")
	}

	functions, err := cmd.List(target, f.kubeClient)
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
	target, err := GetTokenData(s.Token)
	if err != nil {
		return nil, err
	}

	deployConfig := &cmd.DeployHandlerConfig{
		EnableHttpProbe:   f.conf.EnableHttpProbe,
		ImagePullPolicy:   f.conf.ImagePullPolicy,
		FunctionNamespace: f.conf.FunctionNamespace,
		FxWatcherPort:     f.conf.FxWatcherPort,
		SecretMountPath:   f.conf.SecretMountPath,
	}

	err = cmd.Deploy(target, s, f.kubeClient, deployConfig)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}

// grpc handler
// function 삭제
func (f *FxGateway) Delete(c context.Context, s *pb.DeleteFunctionRequest) (*pb.Message, error) {

	target, err := GetTokenData(s.Token) // ?
	if err != nil {
		return nil, err
	}

	if target.Grade == "user" {
		return nil, errors.New("This is an unauthorized API.")
	}

	fnList, err := cmd.List(target, f.kubeClient)
	if err != nil {
		return nil, err
	}

	IsAccessble := IsAuth(fnList, s.FunctionName)
	if !IsAccessble {
		return &pb.Message{Msg: "This is an unauthorized function call."}, errors.New("did not delete: not found function")
	}

	err = cmd.Delete(s.FunctionName, target.Scope, f.kubeClient)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Msg: "OK"}, nil
}

// grpc handler
// function 업데이트
func (f *FxGateway) Update(c context.Context, s *pb.CreateFunctionRequest) (*pb.Message, error) {

	target, err := GetTokenData(s.Token)
	if err != nil {
		return nil, err
	}

	fnList, err := cmd.List(target, f.kubeClient)
	if err != nil {
		return nil, err
	}

	IsAccessble := IsAuth(fnList, s.Service)
	if !IsAccessble {
		return &pb.Message{Msg: "This is an unauthorized function call."}, errors.New("did not update: not found function")
	}

	err = cmd.Update(target.Scope, s, f.kubeClient, f.conf.SecretMountPath)
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

	target, err := GetTokenData(s.Token)
	if err != nil {
		return nil, err
	}

	err = cmd.ReplicaUpdate(target.Scope, s, f.kubeClient)
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

func GetTokenData(token string) (cmd.Access_info, error) {

	// cmd list 랑 service Name 가져와서 in 하고 있으면 true 없으면 false return 하는 로직
	resp, err := http.Get(fmt.Sprintf("%s/verify?access_token=%s", authServerURL, token))
	if err != nil {
		return cmd.Access_info{}, err
	}
	defer resp.Body.Close()

	var target cmd.Access_info
	tokenData, _ := ioutil.ReadAll(resp.Body)

	jsonErr := json.Unmarshal(tokenData, &target)
	if jsonErr != nil {
		return cmd.Access_info{}, fmt.Errorf("cannot parse result from OpenFx on URL: %s\n%s", tokenData, jsonErr.Error())
	}

	return target, nil
}
