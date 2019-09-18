package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/keti-openfx/openfx/cmd"
	"github.com/keti-openfx/openfx/metrics"
	"github.com/keti-openfx/openfx/pb"
	"github.com/keti-openfx/openfx/pkg/ui/data/swagger"
	assetfs "github.com/philips/go-bindata-assetfs"
	"google.golang.org/grpc"
)

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	log.Printf("preflight request for %s", r.URL.Path)
}

func prepareHTTP(ctx context.Context, serverName string, functionNamespace string, fxWatcherPort int, timeout, readtimeout, writetimeout, idletimeout time.Duration) (*http.Server, error) {
	// HTTP router
	router := http.NewServeMux()
	router.Handle("/metrics", metrics.PrometheusHandler())
	router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, strings.NewReader(pb.Swagger))
	})
	serveSwagger(router)
	router.HandleFunc("/function/", makeHttpInvoke(functionNamespace, fxWatcherPort, timeout))

	//// initialize grpc-gateway
	// gRPC dialup options
	opts := []grpc.DialOption{
		grpc.WithTimeout(10 * time.Second),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}

	// gRPC dialup options
	// gRPC server에 대한 클라이언트 연결
	conn, err := grpc.Dial(serverName, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return nil, err
	}

	// changes json serializer to include empty fields with default values
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)

	// Register Gateway endpoints
	// FxGateway 서비스에 대한 http 핸들러를 mux에 등록
	// 핸들러는 conn을 통해 grpc 엔드 포인트로 요청 전달
	err = pb.RegisterFxGatewayHandler(ctx, gwMux, conn)
	if err != nil {
		return nil, err
	}
	////

	router.Handle("/", gwMux)

	// Return HTTP Server instance
	return &http.Server{
		Addr:         serverName,
		Handler:      allowCORS(router),
		ReadTimeout:  readtimeout,
		WriteTimeout: writetimeout,
		IdleTimeout:  idletimeout,
	}, nil
}

func makeHttpInvoke(functionNamespace string, fxWatcherPort int, timeout time.Duration) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Body != nil {
			defer r.Body.Close()
		}

		switch r.Method {
		case http.MethodGet,
			http.MethodPost:

			param := r.URL.Path[10:]
			idx := strings.Index(param, "/")
			var serviceName string
			if idx == -1 {
				serviceName = param
			} else if idx == len(param)-1 {
				serviceName = param[:idx]
			}
			//static := param[idx+1:]

			validName := regexp.MustCompile(`^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`)
			if matched := validName.MatchString(serviceName); !matched {
				buf := bytes.NewBufferString("Must be a valid function name: " + serviceName + "\n")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(buf.Bytes()))
				return
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				buf := bytes.NewBufferString("Error reading body: " + err.Error() + "\n")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(buf.Bytes()))
				return
			}
			output, err := cmd.Invoke(serviceName, functionNamespace, fxWatcherPort, body, timeout)
			if err != nil {
				log.Println(err.Error())
				buf := bytes.NewBufferString("Can't reach service: " + serviceName + "\n")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(buf.Bytes()))
				return
			}

			w.WriteHeader(http.StatusOK)

			w.Write([]byte(output))
			w.Write([]byte("\n"))
		}
	}
}

//func serveSwagger(h *mux.Router) {
func serveSwagger(h *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	h.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
