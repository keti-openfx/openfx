package cmd 

import (
	"fmt"
	"log"

	"github.com/keti-openfx/openfx/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// updates desired count of replicas
func ReplicaUpdate(functionNamespace string, req *pb.ScaleServiceRequest, clientset *kubernetes.Clientset) error {
	log.Println("Update replicas")

	functionName := req.ServiceName

	options := metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
	}
	deployment, err := clientset.ExtensionsV1beta1().Deployments(functionNamespace).Get(functionName, options)

	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, err.Error())
	}

	var replicas int32
	replicas = int32(req.Replicas)
	deployment.Spec.Replicas = &replicas
	_, err = clientset.ExtensionsV1beta1().Deployments(functionNamespace).Update(deployment)

	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil

}

// reads the amount of replicas for a deployment
func GetMeta(functionName string, functionNamespace string, clientset *kubernetes.Clientset) (*pb.Function, error) {
	function, err := getService(functionNamespace, functionName, clientset)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if function == nil {
		return nil, status.Error(codes.NotFound, "function is not exist")
	}

	return function, nil
}

func getService(functionNamespace string, functionName string, clientset *kubernetes.Clientset) (*pb.Function, error) {

	getOpts := metav1.GetOptions{}

	item, err := clientset.ExtensionsV1beta1().Deployments(functionNamespace).Get(functionName, getOpts)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	if item != nil {

		function := readFunction(*item)
		if function != nil {
			return function, nil
		}
	}

	return nil, fmt.Errorf("function: %s not found", functionName)
}
