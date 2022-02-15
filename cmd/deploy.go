package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/keti-openfx/openfx/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"

	//v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"

	v2beta1 "k8s.io/api/autoscaling/v2beta1"
)

/* initialReplicasCount how many replicas to start of creating for a function */
const initialReplicasCount = 1

/* initialCpuUtilization limit of CPU Utilization per pod */
const initialCpuUtilization = 80

type DeployHandlerConfig struct {
	FunctionNamespace string
	EnableHttpProbe   bool
	ImagePullPolicy   string
	FxWatcherPort     int
	FxMeshPort        int
	SecretMountPath   string
}

/* ValidateDeployRequest validates that the service name is valid for Kubernetes */
func ValidateServiceName(service string) error {
	/* Regex for RFC-1123 validation:
	 *	k8s.io/kubernetes/pkg/util/validation/validation.go */
	var validDNS = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	matched := validDNS.MatchString(service)
	if matched {
		return nil
	}

	return fmt.Errorf("(%s) must be a valid DNS entry for service name", service)
}

func Deploy(req *pb.CreateFunctionRequest, clientset *kubernetes.Clientset, config *DeployHandlerConfig) error {
	if err := ValidateServiceName(req.Service); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	log.Printf("Deploying... \ndeploy handler config:%+v\n create function request:%+v\n", config, req)

	existingSecrets, err := getSecrets(clientset, config.FunctionNamespace, req.Secrets)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := context.Background()
	createOpts := metav1.CreateOptions{}
	/*
		persistentVolume := clientset.CoreV1().PersistentVolumes()
		persistentVolumeSpec := makePersistentVolumeSpec(req)
		_, err_pv := persistentVolume.Create(ctx, persistentVolumeSpec, createOpts)

		if err_pv != nil {
			if k8sErrors.IsAlreadyExists(err) {
				return status.Error(codes.AlreadyExists, err.Error())
			}
			return status.Error(codes.Internal, err.Error())
		}

		persistentVolumeClaim := clientset.CoreV1().PersistentVolumeClaims(config.FunctionNamespace)
		persistentVolumeClaimSpec := makePersistentVolumeClaimSpec(req)
		_, err = persistentVolumeClaim.Create(ctx, persistentVolumeClaimSpec, createOpts)

		if err != nil {
			if k8sErrors.IsAlreadyExists(err) {
				return status.Error(codes.AlreadyExists, err.Error())
			}
			return status.Error(codes.Internal, err.Error())
		}
	*/
	deploymentSpec, err := makeDeploymentSpec(req, existingSecrets, config)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	deploy := clientset.AppsV1().Deployments(config.FunctionNamespace)
	_, err = deploy.Create(ctx, deploymentSpec, createOpts)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return status.Error(codes.AlreadyExists, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	log.Println("Created deployment - " + req.Service)

	service := clientset.CoreV1().Services(config.FunctionNamespace)
	serviceSpec := makeServiceSpec(req, config.FxWatcherPort, config.FxMeshPort)
	_, err = service.Create(ctx, serviceSpec, createOpts)

	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return status.Error(codes.AlreadyExists, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	log.Println("Created service - " + req.Service)

	hpaSpec, err := makeAutoscaleSpec(req)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	deployHPA := clientset.AutoscalingV2beta1().HorizontalPodAutoscalers(config.FunctionNamespace)
	_, err = deployHPA.Create(ctx, hpaSpec, createOpts)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return status.Error(codes.AlreadyExists, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	log.Println("Created HPA deployment - " + req.Service)

	return nil
}

func makeDeploymentSpec(req *pb.CreateFunctionRequest, existingSecrets map[string]*apiv1.Secret, config *DeployHandlerConfig) (*v1.Deployment, error) {
	envVars := buildEnvVars(req)
	path := filepath.Join(os.TempDir(), ".lock")
	probe := &apiv1.Probe{
		Handler: apiv1.Handler{
			Exec: &apiv1.ExecAction{
				Command: []string{"cat", path},
			},
		},
		InitialDelaySeconds: 3,
		TimeoutSeconds:      1,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	if config.EnableHttpProbe {
		probe = &apiv1.Probe{
			Handler: apiv1.Handler{
				HTTPGet: &apiv1.HTTPGetAction{
					Path: "/health",
					Port: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(config.FxWatcherPort),
					},
				},
			},
			InitialDelaySeconds: 3,
			TimeoutSeconds:      1,
			PeriodSeconds:       10,
			SuccessThreshold:    1,
			FailureThreshold:    3,
		}

	}

	initialReplicas := int32p(initialReplicasCount)
	labels := map[string]string{
		"kubesphere_openfx_fn_system": "user_fn",
	}

	if req.Labels != nil {
		if min := getMinReplicaCount(req.Labels); min != nil {
			initialReplicas = min
		}
		for k, v := range req.Labels {
			labels[k] = v
		}
	}

	nodeSelector := createSelector(req.Constraints)

	resources, resourceErr := createResources(req)

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

	deploymentSpec := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Service,
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"kubesphere_openfx_fn_system": "user_fn",
				},
			},
			Replicas: initialReplicas,
			Strategy: v1.DeploymentStrategy{
				Type: v1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &v1.RollingUpdateDeployment{
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
					Name:        req.Service,
					Labels:      labels,
					Annotations: buildAnnotations(req),
				},
				Spec: apiv1.PodSpec{
					NodeSelector: nodeSelector,
					//NodeName: "gpu01",
					Containers: []apiv1.Container{
						{
							Name:  req.Service,
							Image: req.Image,
							Ports: []apiv1.ContainerPort{
								{ContainerPort: int32(config.FxWatcherPort), Protocol: apiv1.ProtocolTCP},
								{ContainerPort: int32(config.FxMeshPort), Protocol: apiv1.ProtocolTCP},
							},
							Env:       envVars,
							Resources: *resources,
							/*
								VolumeMounts: []apiv1.VolumeMount{
									{
										Name: req.Service,
										MountPath: "/data",
									},
								},
							*/
							ImagePullPolicy: imagePullPolicy,
							LivenessProbe:   probe,
							ReadinessProbe:  probe,
						},
					},
					/*
						Volumes: []apiv1.Volume{
							{
								Name: req.Service,
								VolumeSource: apiv1.VolumeSource{
									PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
										ClaimName: req.Service,
									},
								},
							},
						},
					*/
					RestartPolicy: apiv1.RestartPolicyAlways,
					DNSPolicy:     apiv1.DNSClusterFirst,
				},
			},
		},
	}

	if err := UpdateSecrets(req, deploymentSpec, existingSecrets, config.SecretMountPath); err != nil {
		return nil, err
	}

	return deploymentSpec, nil
}

func makeServiceSpec(req *pb.CreateFunctionRequest, fxWatcherPort int, fxMeshPort int) *apiv1.Service {
	serviceSpec := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Service,
			Annotations: buildAnnotations(req),
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"kubesphere_openfx_fn_system": "user_fn",
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Name:     "fxwatcher",
					Port:     int32(fxWatcherPort),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(fxWatcherPort),
					},
				},
				{
					Protocol: apiv1.ProtocolTCP,
					Name:     "fxmesh",
					Port:     int32(fxMeshPort),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(fxMeshPort),
					},
				},
			},
		},
	}
	return serviceSpec
}

func makeAutoscaleSpec(req *pb.CreateFunctionRequest) (*v2beta1.HorizontalPodAutoscaler, error) {
	hpaSpec := &v2beta1.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2beta1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: req.Service,
		},

		Spec: v2beta1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2beta1.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       req.Service,
				APIVersion: "apps/v1",
			},
			MinReplicas: int32p(req.MinReplicas),
			MaxReplicas: req.MaxReplicas,
			Metrics: []v2beta1.MetricSpec{
				{
					Type: v2beta1.ResourceMetricSourceType,
					Resource: &v2beta1.ResourceMetricSource{
						Name:                     apiv1.ResourceCPU,
						TargetAverageUtilization: int32p(initialCpuUtilization),
					},
				},

				{
					Type: v2beta1.ResourceMetricSourceType,
					Resource: &v2beta1.ResourceMetricSource{
						Name:               apiv1.ResourceMemory,
						TargetAverageValue: resource.NewQuantity(200*1024*1024, resource.BinarySI),
					},
				},
			},
		},
	}
	return hpaSpec, nil
}

/*
func makePersistentVolumeSpec(req *pb.CreateFunctionRequest) *apiv1.PersistentVolume {
	pvSpec := &apiv1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind: "PersistentVolume",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: req.Service,
		},

		Spec: apiv1.PersistentVolumeSpec{
			Capacity: apiv1.ResourceList{
				apiv1.ResourceStorage: resource.MustParse("10Gi"),
			},
			AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
			PersistentVolumeReclaimPolicy: apiv1.PersistentVolumeReclaimDelete,
			StorageClassName: "local-storage",
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				Local: &apiv1.LocalVolumeSource{
					Path: "/data",
				},
			},
			NodeAffinity: &apiv1.VolumeNodeAffinity{
				Required: &apiv1.NodeSelector{
					NodeSelectorTerms: []apiv1.NodeSelectorTerm{
						{
							MatchExpressions: []apiv1.NodeSelectorRequirement{
								{
									Key: "type",
									Operator: apiv1.NodeSelectorOpIn,
									Values: []string{"gpunode"},
								},
							},
						},
					},
				},
			},
		},
	}
	return pvSpec
}

func makePersistentVolumeClaimSpec(req *pb.CreateFunctionRequest) *apiv1.PersistentVolumeClaim {
	resources := apiv1.ResourceRequirements{
		Requests: apiv1.ResourceList{},
	}

	resources.Requests[apiv1.ResourceStorage] = resource.MustParse("2Gi")

	var storageClassNamePointer *string
	var storageClassName string
	storageClassName = "local-storage"
	storageClassNamePointer = &storageClassName

	pvClaimSpec := &apiv1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind: "PersistentVolumeClaim",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: req.Service,
		},

		Spec: apiv1.PersistentVolumeClaimSpec{
			VolumeName: req.Service,
			AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
			Resources: resources,
			StorageClassName: storageClassNamePointer,
		},
	}
	return pvClaimSpec
}
*/
func buildAnnotations(request *pb.CreateFunctionRequest) map[string]string {
	var annotations map[string]string
	if request.Annotations != nil {
		annotations = request.Annotations
	} else {
		annotations = map[string]string{}
	}

	//annotations["prometheus.io.scrape"] = "false"
	return annotations
}

func buildEnvVars(req *pb.CreateFunctionRequest) []apiv1.EnvVar {
	envVars := []apiv1.EnvVar{}

	for k, v := range req.EnvVars {
		envVars = append(envVars, apiv1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	return envVars
}

func int32p(i int32) *int32 {
	return &i
}

func createSelector(constraints []string) map[string]string {
	selector := make(map[string]string)

	log.Println(constraints)
	if len(constraints) > 0 {
		for _, constraint := range constraints {
			parts := strings.Split(constraint, "=")

			if len(parts) == 2 {
				selector[parts[0]] = parts[1]
			}
		}
	}

	return selector
}

func createResources(req *pb.CreateFunctionRequest) (*apiv1.ResourceRequirements, error) {
	resources := &apiv1.ResourceRequirements{
		Limits:   apiv1.ResourceList{},
		Requests: apiv1.ResourceList{},
	}

	if req.Limits != nil {
		// Set Memory limits
		if len(req.Limits.Memory) > 0 {
			qty, err := resource.ParseQuantity(req.Limits.Memory)
			if err != nil {
				return resources, err
			}
			resources.Limits[apiv1.ResourceMemory] = qty
		}
		// Set CPU limits
		if req.Limits != nil && len(req.Limits.CPU) > 0 {
			qty, err := resource.ParseQuantity(req.Limits.CPU)
			if err != nil {
				return resources, err
			}
			resources.Limits[apiv1.ResourceCPU] = qty
		}
		// Set Gpu limits
		if req.Limits != nil && len(req.Limits.GPU) > 0 {
			qty, err := resource.ParseQuantity(req.Limits.GPU)
			if err != nil {
				return resources, err
			}
			resources.Limits["nvidia.com/gpu"] = qty
		}

	}

	if req.Requests != nil {
		if len(req.Requests.Memory) > 0 {
			qty, err := resource.ParseQuantity(req.Requests.Memory)
			if err != nil {
				return resources, err
			}
			resources.Requests[apiv1.ResourceMemory] = qty
		}

		if len(req.Requests.CPU) > 0 {
			qty, err := resource.ParseQuantity(req.Requests.CPU)
			if err != nil {
				return resources, err
			}
			resources.Requests[apiv1.ResourceCPU] = qty
		}
	}

	return resources, nil
}

func getMinReplicaCount(labels map[string]string) *int32 {
	if value, exists := labels["scale_min"]; exists {
		minReplicas, err := strconv.Atoi(value)
		if err == nil && minReplicas > 0 {
			return int32p(int32(minReplicas))
		}

		log.Println(err)
	}

	return nil
}
