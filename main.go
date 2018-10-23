package main

import (
	"log"

	"github.com/keti-openfx/openfx-gateway/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {

	conf, err := rest.InClusterConfig()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("kubernetest host: %s\n", conf.Host)

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Panic(err)
	}

	c := config.NewFxGatewayConfig()

	s := NewFxGateway(c, clientset)
	log.Printf("[fxgateway] service start\n")
	s.Start()
}
