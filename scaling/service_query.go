// Copyright (c) OpenFaaS Author(s). All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package scaling

import (
	"log"
	"time"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


// ServiceQuery provides interface for replica querying/setting
type ServiceQuery interface {
	GetReplicas(functionName string, namespace string, clientset *kubernetes.Clientset) (response ServiceQueryResponse, err error)
	SetReplicas(functionName string, Namespace string, Replicas int, clientset *kubernetes.Clientset) error
}

// ServiceQueryResponse response from querying a function status
type ServiceQueryResponse struct {
	Replicas          uint64
	MaxReplicas       uint64
	MinReplicas       uint64
	ScalingFactor     uint64
	AvailableReplicas uint64
	Annotations       map[string]string
}


// setReplicas 
func SetReplicas(functionName string, Namespace string, Replicas int, clientset *kubernetes.Clientset) error {

	log.Println("Update replicas")

	options := metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
	}

	deployment, err := clientset.Extensions().Deployments(Namespace).Get(functionName, options)

	if err != nil {
		log.Println(err)
		return err
	}

	oldReplicas := *deployment.Spec.Replicas
	replicas := int32(Replicas) 

	log.Printf("Set replicas - %s %s, %d/%d\n", functionName, Namespace, replicas, oldReplicas)

	deployment.Spec.Replicas = &replicas

	_, err = clientset.Extensions().Deployments(Namespace).Update(deployment)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// getReplicas 
func GetReplicas(functionName string, nameSpace string, clientset *kubernetes.Clientset) (ServiceQueryResponse, error) {
	start := time.Now()
	log.Printf("GetReplicas [%s.%s] took: %fs", functionName, nameSpace, time.Since(start).Seconds())

	options := metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
	}

	dep, err := clientset.Extensions().Deployments(nameSpace).Get(functionName, options)
	if err != nil {
		return ServiceQueryResponse{}, err
	}

	minReplicas := uint64(DefaultMinReplicas)
	maxReplicas := uint64(DefaultMaxReplicas)
	scalingFactor := uint64(DefaultScalingFactor)

	return ServiceQueryResponse{
		Replicas:          uint64(dep.Status.Replicas),
		MaxReplicas:       maxReplicas,
		MinReplicas:       minReplicas,
		ScalingFactor:     scalingFactor,
		AvailableReplicas: uint64(dep.Status.AvailableReplicas),
		Annotations:       dep.GetObjectMeta().GetAnnotations()}, nil
}


