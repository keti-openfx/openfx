package cmd

import (
	"fmt"
	"log"

	"github.com/keti-openfx/openfx/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// IDE deployment
func CreateIDE(req *pb.LoginRequest, clientset *kubernetes.Clientset, config *DeployHandlerConfig) (string, error) {

	persistentVolume := clientset.Core().PersistentVolumes()
	persistentVolumeSpec := makePersistentVolumeSpec(req)
	_, err := persistentVolume.Create(persistentVolumeSpec)

	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return "", status.Error(codes.AlreadyExists, err.Error())
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	persistentVolumeClaim := clientset.Core().PersistentVolumeClaims("openfx-ide")
	persistentVolumeClaimSpec := makePersistentVolumeClaimSpec(req)
	_, err = persistentVolumeClaim.Create(persistentVolumeClaimSpec)

	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return "", status.Error(codes.AlreadyExists, err.Error())
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	existingSecrets, err := getSecrets(clientset, "openfx-ide", []string{"regcred"})
	if err != nil {
		return "", status.Error(codes.InvalidArgument, err.Error())
	}

	deploymentSpec, err := makeIDE_DeploySpec(req, existingSecrets, config)
	if err != nil {
		return "", status.Error(codes.InvalidArgument, err.Error())
	}

	deploy := clientset.Extensions().Deployments("openfx-ide")
	_, err = deploy.Create(deploymentSpec)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return "", status.Error(codes.AlreadyExists, err.Error())
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	service := clientset.Core().Services("openfx-ide")
	serviceSpec := makeIDEServiceSpec(req)
	_, err = service.Create(serviceSpec)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return "", status.Error(codes.AlreadyExists, err.Error())
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	ingresss := clientset.Extensions().Ingresses("openfx-ide")
	ingressSpec := makeIDEIngressSpec(req)
	ingressObj, err := ingresss.Create(ingressSpec)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return "", status.Error(codes.AlreadyExists, err.Error())
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	log.Println("[Succcessfully] Created IDE Env - " + req.Member)

	return ingressObj.Spec.Rules[0].Host + ":32056/" + req.Member + "/", nil
}

func makeIDE_DeploySpec(req *pb.LoginRequest, existingSecrets map[string]*apiv1.Secret, config *DeployHandlerConfig) (*v1beta1.Deployment, error) {
	initialReplicas := int32p(initialReplicasCount)
	labels := map[string]string{
		"openfx-ide": req.Member,
	}

	resources, resourceErr := createResourcesIDE()
	if resourceErr != nil {
		return nil, resourceErr
	}

	var imagePullPolicy apiv1.PullPolicy
	switch config.ImagePullPolicy {
	case "Never":
		imagePullPolicy = apiv1.PullNever
	case "IfNotPresent":
		imagePullPolicy = apiv1.PullIfNotPresent
	default:
		imagePullPolicy = apiv1.PullAlways
	}

	deploymentSpec := &v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Member,
		},
		Spec: v1beta1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"openfx-ide": req.Member,
				},
			},
			Replicas: initialReplicas,
			Strategy: v1beta1.DeploymentStrategy{
				Type: v1beta1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &v1beta1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(1),
					},
				},
			},
			RevisionHistoryLimit: int32p(10),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        req.Member,
					Labels:      labels,
					Annotations: map[string]string{},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  req.Member,
							Image: "10.0.0.91:5000/theia-openfx-ide:0.0.2",
							Ports: []apiv1.ContainerPort{
								{ContainerPort: 3000, Protocol: v1.ProtocolTCP},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      req.Member,
									MountPath: "/app/project",
								},
							},
							ImagePullPolicy: imagePullPolicy,
							Resources:       *resources,
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: req.Member,
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: req.Member,
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyAlways,
					DNSPolicy:     v1.DNSClusterFirst,
				},
			},
		},
	}
	if err := UpdateSecrets_ide(req, deploymentSpec, existingSecrets, config.SecretMountPath); err != nil {
		return nil, err
	}

	return deploymentSpec, nil
}

func makeIDEServiceSpec(req *pb.LoginRequest) *v1.Service {
	serviceSpec := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Member,
			Annotations: map[string]string{},
		},
		Spec: v1.ServiceSpec{
			Type: "ClusterIP",
			Selector: map[string]string{
				"openfx-ide": req.Member,
			},
			Ports: []v1.ServicePort{
				{
					Protocol: v1.ProtocolTCP,
					Port:     4000,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 3000,
					},
				},
			},
		},
	}
	return serviceSpec
}

func createResourcesIDE() (*apiv1.ResourceRequirements, error) {
	resources := &apiv1.ResourceRequirements{
		Limits:   apiv1.ResourceList{},
		Requests: apiv1.ResourceList{},
	}

	/* 	resources.Limits[apiv1.ResourceMemory], _ = resource.ParseQuantity("1Gi")

	   	resources.Limits[apiv1.ResourceCPU], _ = resource.ParseQuantity("1")

	   	resources.Requests[apiv1.ResourceMemory], _ = resource.ParseQuantity("1Gi")

	   	resources.Requests[apiv1.ResourceCPU], _ = resource.ParseQuantity("1")
	*/
	return resources, nil
}

// UpdateSecrets will update the Deployment spec to include secrets that have beenb deployed
// in the kubernetes cluster.  For each requested secret, we inspect the type and add it to the
// deployment spec as appropriat: secrets with type `SecretTypeDockercfg/SecretTypeDockerjson`
// are added as ImagePullSecrets all other secrets are mounted as files in the deployments containers.
func UpdateSecrets_ide(req *pb.LoginRequest, deployment *v1beta1.Deployment, existingSecrets map[string]*apiv1.Secret, secretsMountPath string) error {
	// Add / reference pre-existing secrets within Kubernetes
	secretVolumeProjections := []apiv1.VolumeProjection{}

	secrets := []string{"regcred"}

	for _, secretName := range secrets {
		deployedSecret, ok := existingSecrets[secretName]
		if !ok {
			return fmt.Errorf("Required secret '%s' was not found in the cluster", secretName)
		}

		switch deployedSecret.Type {

		case apiv1.SecretTypeDockercfg,
			apiv1.SecretTypeDockerConfigJson:

			deployment.Spec.Template.Spec.ImagePullSecrets = append(
				deployment.Spec.Template.Spec.ImagePullSecrets,
				apiv1.LocalObjectReference{
					Name: secretName,
				},
			)

			break

		default:

			projectedPaths := []apiv1.KeyToPath{}
			for secretKey := range deployedSecret.Data {
				projectedPaths = append(projectedPaths, apiv1.KeyToPath{Key: secretKey, Path: secretKey})
			}

			projection := &apiv1.SecretProjection{Items: projectedPaths}
			projection.Name = secretName
			secretProjection := apiv1.VolumeProjection{
				Secret: projection,
			}
			secretVolumeProjections = append(secretVolumeProjections, secretProjection)

			break
		}
	}

	volumeName := fmt.Sprintf("%s-projected-secrets", req.Member)
	projectedSecrets := apiv1.Volume{
		Name: volumeName,
		VolumeSource: apiv1.VolumeSource{
			Projected: &apiv1.ProjectedVolumeSource{
				Sources: secretVolumeProjections,
			},
		},
	}

	// remove the existing secrets volume, if we can find it. The update volume will be
	// added below
	existingVolumes := removeVolume(volumeName, deployment.Spec.Template.Spec.Volumes)
	deployment.Spec.Template.Spec.Volumes = existingVolumes
	if len(secretVolumeProjections) > 0 {
		deployment.Spec.Template.Spec.Volumes = append(existingVolumes, projectedSecrets)
	}

	// add mount secret as a file
	updatedContainers := []apiv1.Container{}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		mount := apiv1.VolumeMount{
			Name:      volumeName,
			ReadOnly:  true,
			MountPath: secretsMountPath,
		}

		// remove the existing secrets volume mount, if we can find it. We update it later.
		container.VolumeMounts = removeVolumeMount(volumeName, container.VolumeMounts)
		if len(secretVolumeProjections) > 0 {
			container.VolumeMounts = append(container.VolumeMounts, mount)
		}

		updatedContainers = append(updatedContainers, container)
	}

	deployment.Spec.Template.Spec.Containers = updatedContainers

	return nil
}

func makePersistentVolumeSpec(req *pb.LoginRequest) *v1.PersistentVolume {
	resources := apiv1.ResourceList{}
	resources[v1.ResourceStorage] = resource.MustParse("5Gi")

	pvSpec := &v1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: req.Member,
		},

		Spec: v1.PersistentVolumeSpec{
			Capacity: resources,
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/mnt/data/" + req.Member,
				},
			},
			AccessModes:                   []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			StorageClassName:              "standard",
			PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimRetain,
		},
	}
	return pvSpec
}

func makePersistentVolumeClaimSpec(req *pb.LoginRequest) *v1.PersistentVolumeClaim {
	resources := apiv1.ResourceRequirements{
		Requests: apiv1.ResourceList{},
	}

	resources.Requests[v1.ResourceStorage] = resource.MustParse("5Gi")

	var storageClassNamePointer *string
	var storageClassName string
	storageClassName = "standard"
	storageClassNamePointer = &storageClassName

	pvClaimSpec := &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: req.Member,
		},

		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Resources:        resources,
			VolumeName:       req.Member,
			StorageClassName: storageClassNamePointer,
		},
	}

	return pvClaimSpec
}

func makeIDEIngressSpec(req *pb.LoginRequest) *v1beta1.Ingress {

	Hosts := "keti.asuscomm.com"
	IngressPath := "/" + req.Member + "/(.*)"

	IngressSpec := &v1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Member,
			Namespace:   "openfx-ide",
			Annotations: map[string]string{"nginx.ingress.kubernetes.io/rewrite-target": "/$1"},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: Hosts,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path:    IngressPath,
									Backend: v1beta1.IngressBackend{ServiceName: req.Member, ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 4000}},
								},
							},
						},
					},
				},
			},
		},
	}
	return IngressSpec
}
