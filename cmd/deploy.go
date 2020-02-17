package cmd

import (
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
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
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

	deploymentSpec, err := makeDeploymentSpec(req, existingSecrets, config)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	deploy := clientset.Extensions().Deployments(config.FunctionNamespace)
	_, err = deploy.Create(deploymentSpec)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return status.Error(codes.AlreadyExists, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	log.Println("Created deployment - " + req.Service)

	service := clientset.Core().Services(config.FunctionNamespace)
	serviceSpec := makeServiceSpec(req, config.FxWatcherPort)
	_, err = service.Create(serviceSpec)

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
	_, err = deployHPA.Create(hpaSpec)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return status.Error(codes.AlreadyExists, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

	log.Println("Created HPA deployment - " + req.Service)

	return nil
}

func makeDeploymentSpec(req *pb.CreateFunctionRequest, existingSecrets map[string]*apiv1.Secret, config *DeployHandlerConfig) (*v1beta1.Deployment, error) {
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
		"openfx_fn": req.Service,
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

	deploymentSpec := &v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Service,
		},
		Spec: v1beta1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"openfx_fn": req.Service,
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
					Name:        req.Service,
					Labels:      labels,
					Annotations: buildAnnotations(req),
				},
				Spec: apiv1.PodSpec{
					NodeSelector: nodeSelector,
					Containers: []apiv1.Container{
						{
							Name:  req.Service,
							Image: req.Image,
							Ports: []apiv1.ContainerPort{
								{ContainerPort: int32(config.FxWatcherPort), Protocol: v1.ProtocolTCP},
							},
							Env:             envVars,
							Resources:       *resources,
							ImagePullPolicy: imagePullPolicy,
							LivenessProbe:   probe,
							ReadinessProbe:  probe,
						},
					},
					RestartPolicy: v1.RestartPolicyAlways,
					DNSPolicy:     v1.DNSClusterFirst,
				},
			},
		},
	}

	if err := UpdateSecrets(req, deploymentSpec, existingSecrets, config.SecretMountPath); err != nil {
		return nil, err
	}

	return deploymentSpec, nil
}

func makeServiceSpec(req *pb.CreateFunctionRequest, fxWatcherPort int) *v1.Service {
	serviceSpec := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Service,
			Annotations: buildAnnotations(req),
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"openfx_fn": req.Service,
			},
			Ports: []v1.ServicePort{
				{
					Protocol: v1.ProtocolTCP,
					Port:     int32(fxWatcherPort),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(fxWatcherPort),
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
				APIVersion: "extensions/v1beta1",
			},
			//MinReplicas: int32p(initialReplicasCount),
			MinReplicas: int32p(req.MinReplicas),
			//MaxReplicas: int32(5),
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

func buildEnvVars(req *pb.CreateFunctionRequest) []v1.EnvVar {
	envVars := []v1.EnvVar{}

	for k, v := range req.EnvVars {
		envVars = append(envVars, v1.EnvVar{
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

		// Set GPU limits
		if req.Limits != nil && len(req.Limits.GPU) > 0 {
			qty, err := resource.ParseQuantity(req.Limits.GPU)
			if err != nil {
				return resources, err
			}
			resources.Limits[apiv1.ResourceNvidiaGPU] = qty
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
