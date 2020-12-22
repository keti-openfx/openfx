package main

import (

	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/keti-openfx/fx-idler/types"
	"github.com/keti-openfx/openfx/metrics"
	"github.com/keti-openfx/fx-idler/api/grpc"
	"github.com/keti-openfx/fx-idler/pb"
)

const (
	scaleLabel     = "openfx_fn"
	defaultNamespace = "openfx-fn"
)

var dryRun bool

var writeDebug bool

func main() {
	config, configErr := types.ReadConfig() 
	if configErr != nil {
		log.Panic(configErr.Error())
		os.Exit(1)
	}

	flag.BoolVar(&dryRun, "dry-run", false, "use dry-run for scaling events")
	flag.Parse()

	if val, ok := os.LookupEnv("write_debug"); ok && (val == "1" || val == "true") {
		writeDebug = true
	}


	log.Printf(`dry_run: %t
gateway_url: %s
inactivity_duration: %s
reconcile_interval: %s
`, dryRun, config.GatewayURL, config.InactivityDuration, config.ReconcileInterval)

	if len(config.GatewayURL) == 0 {
		log.Println("gateway_url (openfx gateway) is required.")
		os.Exit(1)
	}
	// time Loop  
	for {
		reconcile(config)
		time.Sleep(config.ReconcileInterval) 
		log.Printf("\n")
	}
}

//NewHTTPClient returns a new HTTP client
func NewHTTPClient() *http.Client {
	return &http.Client{}
}

func buildMetricsMap(fnList *pb.Functions, config types.Config, namespace string) map[string]float64 {
	query := metrics.NewPrometheusQuery(config.PrometheusHost, config.PrometheusPort, NewHTTPClient())

	duration := fmt.Sprintf("%dm", int(config.InactivityDuration.Minutes()))
	// duration := "5m"
	metricsMap := make(map[string]float64)

	for _, function := range fnList.Functions {
		//Deriving the function name for multiple namespace support

		functionName := fmt.Sprintf("%s", function.Name)


		querySt := url.QueryEscape(fmt.Sprintf(
			`sum(rate(gateway_function_invocation_total{function_name="%s"}[%s])) by (function_name)`,
			functionName,
			duration))

		log.Printf("Query: %s\n", querySt)

		res, err := query.Fetch(querySt)
		if err != nil {
			log.Println(err)
			continue
		}

		if len(res.Data.Result) > 0 || function.InvocationCount == 0 {

			if _, exists := metricsMap[functionName]; !exists {
				metricsMap[functionName] = 0
			}

			for _, v := range res.Data.Result {

				if writeDebug {
					log.Println(v)
				}

				if v.Metric.FunctionName == functionName {
					metricValue := v.Value[1]
					switch metricValue.(type) {
					case string:

						f, strconvErr := strconv.ParseFloat(metricValue.(string), 64)
						if strconvErr != nil {
							log.Printf("Unable to convert value for metric: %s\n", strconvErr)
							continue
						}

						metricsMap[functionName] = metricsMap[functionName] + f
					}
				}
			}
		}
	}
	return metricsMap
}

func reconcile( config types.Config) {
	//First reconcile with openfx-fn namespace for
	reconcileNamespace(config, defaultNamespace) 
}


func reconcileNamespace(config types.Config, namespace string) {
	fnList, err := grpc.List(config.GatewayURL)
	if err != nil {
		log.Println(err)
		return
	}

	metricsMap := buildMetricsMap(fnList, config, namespace)  

	for _, fn := range fnList.Functions{
		//Deriving the function name for multiple namespace support
		functionName := fmt.Sprintf("%s", fn.Name)

/* 		if fn.Labels != nil {
			labels := fn.Labels 
			labelValue := labels[scaleLabel]

			log.Printf("Value: %s %s\n", labels, labelValue)
			if labelValue != "1" && labelValue != "true" {
				if writeDebug {
					log.Printf("Skip: %s due to missing label\n", functionName)
				}
				continue
			}
		} */

		if v, found := metricsMap[functionName]; found {
			if v == float64(0) {
				log.Printf("%s\tidle\n", functionName)
				// GetfunctionInfo -> getMeta fillyo
				val, err := grpc.GetMeta(fn.Name, config.GatewayURL)

				if err != nil {
					log.Println(err.Error())
					continue
				}

				replicaCount := uint64(0)
				if dryRun {
					log.Printf("dry-run: Scaling %s to %d replicas\n", fn.Name, replicaCount)
					continue
				}

				if err == nil && val.AvailableReplicas > 0 {
					_, err = grpc.Scale(config.GatewayURL, fn.Name, namespace, replicaCount)

					if err != nil {
						log.Println(err.Error())
					} else {
						log.Printf("scaled function %s to %d replica(s)\n", fn.Name, replicaCount)
					}
				}
			} else {
				if writeDebug {
					log.Printf("%s\tactive: %f\n", functionName, v)
				}
			}
		}
	}
}
