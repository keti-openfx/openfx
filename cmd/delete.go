package cmd

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Delete(functionName, functionNamespace string, clientset *kubernetes.Clientset) error {

	if err := ValidateServiceName(functionName); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	getOpts := metav1.GetOptions{}

	// This makes sure we don't delete non-labelled deployments
	deployment, findDeployErr := clientset.ExtensionsV1beta1().
		Deployments(functionNamespace).
		Get(functionName, getOpts)

	if findDeployErr != nil {
		if errors.IsNotFound(findDeployErr) {
			return status.Error(codes.NotFound, findDeployErr.Error())
		}
		return status.Error(codes.Internal, findDeployErr.Error())
	}

	if isFunction(deployment) {
		if err := deleteFunction(functionName, functionNamespace, clientset); err != nil {
			return err
		}
	} else {
		return status.Error(codes.Internal, "Not a function: "+functionName)
	}
	return nil
}

func isFunction(deployment *v1beta1.Deployment) bool {
	if deployment != nil {
		if _, found := deployment.Labels["openfx_fn"]; found {
			return true
		}
	}
	return false
}

func deleteFunction(functionName, functionNamespace string, clientset *kubernetes.Clientset) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}

	if deployErr := clientset.ExtensionsV1beta1().
		Deployments(functionNamespace).
		Delete(functionName, opts); deployErr != nil {

		if errors.IsNotFound(deployErr) {
			return status.Error(codes.NotFound, deployErr.Error())

		}
		return status.Error(codes.Internal, deployErr.Error())
	}

	if svcErr := clientset.CoreV1().
		Services(functionNamespace).
		Delete(functionName, opts); svcErr != nil {

		if errors.IsNotFound(svcErr) {
			return status.Error(codes.NotFound, svcErr.Error())
		}
		return status.Error(codes.Internal, svcErr.Error())
	}

	if hpaErr := clientset.AutoscalingV2beta1().
                HorizontalPodAutoscalers(functionNamespace).
                Delete(functionName, opts); hpaErr != nil {

		if errors.IsNotFound(hpaErr) {
			return status.Error(codes.NotFound, hpaErr.Error())
		}
		return status.Error(codes.Internal, hpaErr.Error())
	}

	return nil
}
