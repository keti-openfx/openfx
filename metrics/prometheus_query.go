package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/keti-openfx/openfx/pb"
)

type VectorQueryResponse struct {
	Data struct {
		Result []struct {
			Metric struct {
				Code         string `json:"code"`
				FunctionName string `json:"function_name"`
			}
			Value []interface{} `json:"value"`
		}
	}
}

type PrometheusQueryFetcher interface {
	Fetch(query string) (*VectorQueryResponse, error)
}

// PrometheusQuery represents parameters for querying Prometheus
type PrometheusQuery struct {
	Port   int
	Host   string
	Client *http.Client
}

// NewPrometheusQuery create a NewPrometheusQuery
func NewPrometheusQuery(host string, port int, client *http.Client) PrometheusQuery {
	return PrometheusQuery{
		Client: client,
		Host:   host,
		Port:   port,
	}
}

// Fetch queries aggregated stats
func (q PrometheusQuery) Fetch(query string) (*VectorQueryResponse, error) {

	req, reqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/api/v1/query?query=%s", q.Host, q.Port, query), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	res, getErr := q.Client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	bytesOut, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code from Prometheus want: %d, got: %d, body: %s", http.StatusOK, res.StatusCode, string(bytesOut))
	}

	var values VectorQueryResponse

	unmarshalErr := json.Unmarshal(bytesOut, &values)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("Error unmarshaling result: %s, '%s'", unmarshalErr, string(bytesOut))
	}

	return &values, nil
}

func AddMetricsFunction(fn *pb.Function, prometheusQuery PrometheusQueryFetcher) *pb.Function {
	q := fmt.Sprintf("sum(gateway_function_invocation_total{function_name=~\"%s\", code=~\".*\"}) by (function_name, code)", fn.Name)
	expr := url.QueryEscape(q)
	results, fetchErr := prometheusQuery.Fetch(expr)
	if fetchErr != nil {
		log.Printf("Error querying Prometheus API: %s\n", fetchErr.Error())
		return fn
	}

	fn.InvocationCount = 0
	for _, v := range results.Data.Result {
		if v.Metric.FunctionName == fn.Name {
			metricValue := v.Value[1]
			switch metricValue.(type) {
			case string:
				f, strconvErr := strconv.ParseUint(metricValue.(string), 10, 64)
				if strconvErr != nil {
					log.Printf("Unable to convert value for metric: %s\n", strconvErr)
					continue
				}
				fn.InvocationCount += f
				break
			}
		}
	}

	return fn
}

func AddMetricsFunctions(fns []*pb.Function, prometheusQuery PrometheusQueryFetcher) []*pb.Function {
	expr := url.QueryEscape(`sum(gateway_function_invocation_total{function_name=~".*", code=~".*"}) by (function_name, code)`)
	results, fetchErr := prometheusQuery.Fetch(expr)

	if fetchErr != nil {
		log.Printf("Error querying Prometheus API: %s\n", fetchErr.Error())

		return fns
	}

	return _mixIn(fns, results)
}

func _mixIn(functions []*pb.Function, metrics *VectorQueryResponse) []*pb.Function {
	if functions == nil {
		return nil
	}

	// Ensure values are empty first.
	for _, function := range functions {
		function.InvocationCount = 0
	}

	for _, function := range functions {
		for _, v := range metrics.Data.Result {

			if v.Metric.FunctionName == function.Name {
				metricValue := v.Value[1]
				switch metricValue.(type) {
				case string:
					f, strconvErr := strconv.ParseUint(metricValue.(string), 10, 64)
					if strconvErr != nil {
						log.Printf("Unable to convert value for metric: %s\n", strconvErr)
						continue
					}
					function.InvocationCount += f
					break
				}
			}
		}
	}

	return functions
}
