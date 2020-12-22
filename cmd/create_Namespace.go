package cmd

import (
	"encoding/json"
	"log"

	"github.com/keti-openfx/openfx/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
)

func Create(req *pb.CreateNamespaceRequest, clientset *kubernetes.Clientset) error {

	log.Printf("Creating %v namespace environment according to client registration\n", req.NamespaceName)

	// create namespace
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: req.NamespaceName}}

	_, err := clientset.Core().Namespaces().Create(nsSpec)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return status.Error(codes.AlreadyExists, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	// Get Secret info and create Secret to namespace
	existingSecrets, err := getSecrets(clientset, "openfx", []string{"regcred"})
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	secret_byte, _ := json.Marshal(existingSecrets["regcred"]) // json byte

	secrets := apiv1.Secret{}
	err = json.Unmarshal(secret_byte, &secrets)
	if err != nil {
		panic(err)
	}
	in := []byte(secrets.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"])
	var raw map[string]interface{}
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}

	//result := raw["data"].(map[string]interface{})

	secretspec := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "regcred",
			Namespace: req.NamespaceName,
		},
		Data: map[string][]byte{
			".dockerconfigjson": []byte("{\"auths\":{\"10.0.0.91:5000\":{\"username\":\"test\",\"password\":\"test\",\"auth\":\"dGVzdDp0ZXN0\"}}}"),
		},
		Type: "kubernetes.io/dockerconfigjson",
	}

	_, err = clientset.Core().Secrets(req.NamespaceName).Create(&secretspec)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return status.Error(codes.AlreadyExists, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	log.Printf("Successfully Create %v namespace environment\n", req.NamespaceName)
	return nil
}
