package cmd

import (
	"log"

	"github.com/keti-openfx/openfx/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func List(functionNamespace string, clientset *kubernetes.Clientset) ([]*pb.Function, error) {
	listOpts := metav1.ListOptions{
		LabelSelector: "openfx_fn",
	}
	res, err := clientset.ExtensionsV1beta1().Deployments(functionNamespace).List(listOpts)

	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	functions := []*pb.Function{}
	for _, item := range res.Items {
		function := readFunction(item)
		if function != nil {
			functions = append(functions, function)
		}
	}
	return functions, nil
}

func readFunction(item v1beta1.Deployment) *pb.Function {
	var replicas uint64
	if item.Spec.Replicas != nil {
		replicas = uint64(*item.Spec.Replicas)
	}

	labels := item.Labels
	function := pb.Function{
		Name:              item.Name,
		Replicas:          replicas,
		Image:             item.Spec.Template.Spec.Containers[0].Image,
		AvailableReplicas: uint64(item.Status.AvailableReplicas),
		InvocationCount:   0,
		Labels:            labels,
		Annotations:       item.Spec.Template.Annotations,
	}

	return &function
}
