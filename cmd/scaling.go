package cmd

import (
	"log"
	"strconv"
	"sync"
	"time"

	"k8s.io/client-go/kubernetes"
)

var (
	cache                FunctionCache
	maxPollCount         = uint(1000)
	functionPollInterval = time.Millisecond * 10
)

func init() {
	cache = FunctionCache{
		Cache:  make(map[string]*FunctionMeta),
		Expiry: time.Second * 5,
	}
}

type ServiceReplicas struct {
	Replicas          uint64
	MaxReplicas       uint64
	MinReplicas       uint64
	ScalingFactor     uint64
	AvailableReplicas uint64
}

// FunctionMeta holds the last refresh and any other
// meta-data needed for caching.
type FunctionMeta struct {
	LastRefresh     time.Time
	ServiceReplicas ServiceReplicas
}

// Expired find out whether the cache item has expired with
// the given expiry duration from when it was stored.
func (fm *FunctionMeta) Expired(expiry time.Duration) bool {
	return time.Now().After(fm.LastRefresh.Add(expiry))
}

// FunctionCache provides a cache of Function replica counts
type FunctionCache struct {
	Cache  map[string]*FunctionMeta
	Expiry time.Duration
	Sync   sync.Mutex
}

// Set replica count for functionName
func (fc *FunctionCache) Set(functionName string, serviceReplicas ServiceReplicas) {
	fc.Sync.Lock()
	if _, exists := fc.Cache[functionName]; !exists {
		fc.Cache[functionName] = &FunctionMeta{}
	}

	entry := fc.Cache[functionName]
	entry.LastRefresh = time.Now()
	entry.ServiceReplicas = serviceReplicas

	fc.Sync.Unlock()
}

// Get replica count for functionName
func (fc *FunctionCache) Get(functionName string) (ServiceReplicas, bool) {
	replicas := ServiceReplicas{
		AvailableReplicas: 0,
	}

	hit := false
	fc.Sync.Lock()

	if val, exists := fc.Cache[functionName]; exists {
		replicas = val.ServiceReplicas
		hit = !val.Expired(fc.Expiry)
	}

	fc.Sync.Unlock()
	return replicas, hit
}

func Scaling(functionName string, functionNamespace string, clientset *kubernetes.Clientset) {
	if serviceReplicas, hit := cache.Get(functionName); hit && serviceReplicas.AvailableReplicas > 0 {
		return
	}

	fn, err := GetMeta(functionName, functionNamespace, clientset)
	minReplicas := extractLabelValue(fn.Labels["scale_min"], 1)
	maxReplicas := extractLabelValue(fn.Labels["scale_max"], 20)
	scalingFactor := extractLabelValue(fn.Labels["scale_factor"], 20)
	serviceReplicas := ServiceReplicas{
		MinReplicas:       minReplicas,
		MaxReplicas:       maxReplicas,
		ScalingFactor:     scalingFactor,
		AvailableReplicas: fn.AvailableReplicas,
		Replicas:          fn.Replicas,
	}

	cache.Set(functionName, serviceReplicas)

	if err != nil {
		log.Println(err)
		//return status.Error(codes.Internal, err.Error())
	}

	if serviceReplicas.AvailableReplicas == 0 {
		//set replicas
	}

}

func extractLabelValue(rawLabelValue string, fallback uint64) uint64 {
	if len(rawLabelValue) <= 0 {
		return fallback
	}

	value, err := strconv.Atoi(rawLabelValue)

	if err != nil {
		log.Printf("Provided label value %s should be of type uint", rawLabelValue)
		return fallback
	}

	return uint64(value)
}
