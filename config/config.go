package config

import (
	"time"
)

var FxGatewayVersion string

type FxGatewayConfig struct {
	FunctionNamespace string
	TCPPort           int
	InvokeTimeout     time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ImagePullPolicy   string
	EnableHttpProbe   bool
	FxWatcherPort     int
	FxMeshPort        int
	BasicAuth         bool
	SecretMountPath   string

	PrometheusHost string
	PrometheusPort int
}

func NewFxGatewayConfig(version string) FxGatewayConfig {

	envs := NewEnvs()

	FxGatewayVersion = version

	return FxGatewayConfig{
		FunctionNamespace: envs.GetString("FUNCTION_NAMESPACE", "kubesphere-openfx-fn-system"),
		TCPPort:           envs.GetInt("PORT", 10000),
		InvokeTimeout:     envs.GetDuration("INVOKE_TIMEOUT", time.Second*605),
		ReadTimeout:       envs.GetDuration("READ_TIMEOUT", time.Second*605),
		WriteTimeout:      envs.GetDuration("WRITE_TIMEOUT", time.Second*605),
		IdleTimeout:       envs.GetDuration("IDLE_TIMEOUT", time.Second*120),
		ImagePullPolicy:   envs.GetString("IMAGE_PULL_POLICY", "Always"),
		EnableHttpProbe:   envs.GetBool("FUNCTION_HTTP_PROBE", false),
		FxWatcherPort:     envs.GetInt("FXWATCHER_PORT", 50051),
		FxMeshPort:        envs.GetInt("FXMESH_PORT", 50052),
		BasicAuth:         envs.GetBool("BASIC_AUTH", false),
		SecretMountPath:   envs.GetString("SECRET_MOUNT_PATH", "/etc/openfx"),

		PrometheusHost: envs.GetString("PROMETHEUS_HOST", "prometheus"),
		PrometheusPort: envs.GetInt("PROMETHEUS_PORT", 9090),
	}
}
