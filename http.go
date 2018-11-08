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
	"github.com/keti-openfx/openfx-gateway/metrics"
	"github.com/keti-openfx/openfx-gateway/pb"
	"github.com/keti-openfx/openfx-gateway/pkg/ui/data/swagger"
	"github.com/keti-openfx/openfx-gateway/service"
	assetfs "github.com/philips/go-bindata-assetfs"
	"google.golang.org/grpc"
)

func prepareHTTP(ctx context.Context, serverName string, functionNamespace string, fxWatcherPort int, timeout, readtimeout, writetimeout, idletimeout time.Duration) (*http.Server, error) {
	// HTTP router
	router := http.NewServeMux()
	router.Handle("/metrics", metrics.PrometheusHandler())
	router.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
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
	err = pb.RegisterFxGatewayHandler(ctx, gwMux, conn)
	if err != nil {
		return nil, err
	}
	////

	router.Handle("/", gwMux)

	// Return HTTP Server instance
	return &http.Server{
		Addr:         serverName,
		Handler:      router,
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
			output, err := service.Invoke(serviceName, functionNamespace, fxWatcherPort, body, timeout)
			if err != nil {
				log.Println(err.Error())
				buf := bytes.NewBufferString("Can't reach service: " + serviceName + "\n")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(buf.Bytes()))
				return
			}

			log.Printf("Success Invoke Service....\n")
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
