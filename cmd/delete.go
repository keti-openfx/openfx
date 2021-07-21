package cmd

import (
	"log"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//v1beta1 "k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Delete(functionName, functionNamespace string, clientset *kubernetes.Clientset) error {

	if err := ValidateServiceName(functionName); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := context.Background()
	getOpts := metav1.GetOptions{}

	// This makes sure we don't delete non-labelled deployments
	deployment, findDeployErr := clientset.AppsV1().
		Deployments(functionNamespace).
		Get(ctx, functionName, getOpts)

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

func isFunction(deployment *v1.Deployment) bool {
	if deployment != nil {
		log.Println("Deleted service function ...")
		if _, found := deployment.Spec.Template.ObjectMeta.Labels["openfx_fn"]; found {
			return true
		}
	}
	return false
}

func deleteFunction(functionName, functionNamespace string, clientset *kubernetes.Clientset) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}
	ctx := context.Background()

	if deployErr := clientset.AppsV1().
		Deployments(functionNamespace).
		Delete(ctx, functionName, opts); deployErr != nil {

		if errors.IsNotFound(deployErr) {
			return status.Error(codes.NotFound, deployErr.Error())

		}
		return status.Error(codes.Internal, deployErr.Error())
	}

	if svcErr := clientset.CoreV1().
		Services(functionNamespace).
		Delete(ctx, functionName, opts); svcErr != nil {

		if errors.IsNotFound(svcErr) {
			return status.Error(codes.NotFound, svcErr.Error())
		}
		return status.Error(codes.Internal, svcErr.Error())
	}

	if hpaErr := clientset.AutoscalingV2beta1().
                HorizontalPodAutoscalers(functionNamespace).
                Delete(ctx, functionName, opts); hpaErr != nil {

		if errors.IsNotFound(hpaErr) {
			return status.Error(codes.NotFound, hpaErr.Error())
		}
		return status.Error(codes.Internal, hpaErr.Error())
	}
	
	if pvErr := clientset.CoreV1().
		PersistentVolumes().
		Delete(ctx, functionName, opts); pvErr != nil {
		
		if errors.IsNotFound(pvErr) {
			return status.Error(codes.NotFound, pvErr.Error())
		}
		return status.Error(codes.Internal, pvErr.Error())
	}

	if pvcErr := clientset.CoreV1().
		PersistentVolumeClaims(functionNamespace).
		Delete(ctx, functionName, opts); pvcErr != nil {

		if errors.IsNotFound(pvcErr) {
			return status.Error(codes.NotFound, pvcErr.Error())
		}
		return status.Error(codes.Internal, pvcErr.Error())
	}

	return nil
}
