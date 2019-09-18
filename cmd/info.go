package cmd

import (
	"log"

	"github.com/keti-openfx/openfx/config"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
)

type information struct {
	kubernetes struct {
		version string
	}
	fxGateway struct {
		version string
	}
}

func Info(clientset *kubernetes.Clientset) (string, error) {
	v, err := clientset.ServerVersion()
	if err != nil {
		log.Println(err)
		return "", err
	}

	info := information{}

	info.kubernetes.version = v.GitVersion
	info.fxGateway.version = config.FxGatewayVersion

	gwInfo, err := yaml.Marshal(&info)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(gwInfo), nil
}
