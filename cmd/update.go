package cmd

import (
	"fmt"
	"log"
	"time"
	"context"

	"github.com/keti-openfx/openfx/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Update(functionNamespace string, req *pb.CreateFunctionRequest, clientset *kubernetes.Clientset, secretMountPath string) error {
	if err := ValidateServiceName(req.Service); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	log.Printf("Updating... \ncreate function request:%+v\n", req)

	annotations := buildAnnotations(req)

	if err := updateDeploymentSpec(functionNamespace, clientset, req, annotations, secretMountPath); err != nil {
		return err
	}

	if err := updateService(functionNamespace, clientset, req, annotations); err != nil {
		return err
	}

	if err := updateHPA(functionNamespace, clientset, req, annotations); err != nil {
		return err
	}

	log.Println("Updated service - " + req.Service)

	return nil
}

func updateDeploymentSpec(
	functionNamespace string,
	clientset *kubernetes.Clientset,
	request *pb.CreateFunctionRequest,
	annotations map[string]string,
	secretMountPath string) (err error) {

	ctx := context.Background()
	getOpts := metav1.GetOptions{}
	updateOpts := metav1.UpdateOptions{}
	deployment, findDeployErr := clientset.AppsV1().
		Deployments(functionNamespace).
		Get(ctx, request.Service, getOpts)

	if findDeployErr != nil {
		log.Println(findDeployErr)
		return status.Error(codes.NotFound, findDeployErr.Error())
	}

	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		deployment.Spec.Template.Spec.Containers[0].Image = request.Image

		// Disabling update support to prevent unexpected mutations of deployed functions,
		// since imagePullPolicy is now configurable. This could be reconsidered later depending
		// on desired behavior, but will need to be updated to take config.
		//deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy = v1.PullAlways

		deployment.Spec.Template.Spec.Containers[0].Env = buildEnvVars(request)

		//configureReadOnlyRootFilesystem(request, deployment)

		deployment.Spec.Template.Spec.NodeSelector = createSelector(request.Constraints)

		labels := map[string]string{
			"openfx_fn": request.Service,
			"uid":         fmt.Sprintf("%d", time.Now().Nanosecond()),
		}

		if request.Labels != nil {
			if min := getMinReplicaCount(request.Labels); min != nil {
				deployment.Spec.Replicas = min
			}

			for k, v := range request.Labels {
				labels[k] = v
			}
		}

		deployment.Labels = labels
		deployment.Spec.Template.ObjectMeta.Labels = labels

		deployment.Annotations = annotations
		deployment.Spec.Template.Annotations = annotations
		deployment.Spec.Template.ObjectMeta.Annotations = annotations

		resources, resourceErr := createResources(request)
		if resourceErr != nil {
			log.Println(resourceErr)
			return status.Error(codes.InvalidArgument, resourceErr.Error())
		}

		deployment.Spec.Template.Spec.Containers[0].Resources = *resources

		existingSecrets, err := getSecrets(clientset, functionNamespace, request.Secrets)
		if err != nil {
			log.Println(err)
			return status.Error(codes.InvalidArgument, err.Error())
		}

		err = UpdateSecrets(request, deployment, existingSecrets, secretMountPath)
		if err != nil {
			log.Println(err)
			return status.Error(codes.InvalidArgument, err.Error())
		}
	}

	if _, updateErr := clientset.AppsV1().
		Deployments(functionNamespace).
		Update(ctx, deployment, updateOpts); updateErr != nil {

		log.Println(updateErr)
		return status.Error(codes.Internal, updateErr.Error())
	}

	return nil
}

func updateService(
	functionNamespace string,
	clientset *kubernetes.Clientset,
	request *pb.CreateFunctionRequest,
	annotations map[string]string) (err error) {

	ctx := context.Background()
	getOpts := metav1.GetOptions{}
	updateOpts := metav1.UpdateOptions{}

	service, findServiceErr := clientset.CoreV1().
		Services(functionNamespace).
		Get(ctx, request.Service, getOpts)

	if findServiceErr != nil {
		log.Println(findServiceErr)
		return status.Error(codes.NotFound, findServiceErr.Error())
	}

	service.Annotations = annotations

	if _, updateErr := clientset.CoreV1().
		Services(functionNamespace).
		Update(ctx, service, updateOpts); updateErr != nil {

		log.Println(updateErr)
		return status.Error(codes.Internal, updateErr.Error())
	}

	return nil
}

func updateHPA(
	functionNamespace string,
	clientset *kubernetes.Clientset,
	request *pb.CreateFunctionRequest,
	annotations map[string]string) (err error) {

	ctx := context.Background()
	getOpts := metav1.GetOptions{}
	updateOpts := metav1.UpdateOptions{}
	hpa, findHPAErr := clientset.AutoscalingV2beta1().
		HorizontalPodAutoscalers(functionNamespace).
		Get(ctx, request.Service, getOpts)

	if findHPAErr != nil {
		log.Println(findHPAErr)
		return status.Error(codes.NotFound, findHPAErr.Error())
	}

	hpa.Annotations = annotations
	hpa.Spec.MinReplicas = int32p(request.MinReplicas)
	hpa.Spec.MaxReplicas = request.MaxReplicas

	if _, updateErr := clientset.AutoscalingV2beta1().
		HorizontalPodAutoscalers(functionNamespace).
		Update(ctx, hpa, updateOpts); updateErr != nil {

		log.Println(updateErr)
		return status.Error(codes.Internal, updateErr.Error())
	}

	return nil
}

func updatePVC(
	functionNamespace string,
	clientset *kubernetes.Clientset,
	request *pb.CreateFunctionRequest,
	annotations map[string]string) (err error) {

	ctx := context.Background()
	getOpts := metav1.GetOptions{}
	updateOpts := metav1.UpdateOptions{}
	pvc, findPVCErr := clientset.CoreV1().
		PersistentVolumeClaims(functionNamespace).
		Get(ctx, request.Service, getOpts)

	if findPVCErr != nil {
		log.Println(findPVCErr)
		return status.Error(codes.NotFound, findPVCErr.Error())
	}

	pvc.Annotations = annotations

	if _, updateErr := clientset.CoreV1().
		PersistentVolumeClaims(functionNamespace).
		Update(ctx, pvc, updateOpts); updateErr != nil {

		log.Println(updateErr)
		return status.Error(codes.Internal, updateErr.Error())
	}

	return nil
}
