package config

import (
	"time"
)

const (
	FxGatewayVersion = "0.1.0"
)

type FxGatewayConfig struct {
	FunctionNamespace string
	TCPPort           int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ImagePullPolicy   string
	EnableHttpProbe   bool
	FxWatcherPort     int
	BasicAuth         bool
	SecretMountPath   string
}

func NewFxGatewayConfig() FxGatewayConfig {

	envs := NewEnvs()

	return FxGatewayConfig{
		FunctionNamespace: envs.GetString("FUNCTION_NAMESPACE", "openfx-fn"),
		TCPPort:           envs.GetInt("PORT", 10000),
		ReadTimeout:       envs.GetDuration("READ_TIMEOUT", time.Second*10),
		WriteTimeout:      envs.GetDuration("WRITE_TIMEOUT", time.Second*10),
		IdleTimeout:       envs.GetDuration("IDLE_TIMEOUT", time.Second*10),
		ImagePullPolicy:   envs.GetString("IMAGE_PULL_POLICY", "Always"),
		EnableHttpProbe:   envs.GetBool("FUNCTION_HTTP_PROBE", false),
		FxWatcherPort:     envs.GetInt("FXWATCHER_PORT", 50051),
		BasicAuth:         envs.GetBool("BASIC_AUTH", false),
		SecretMountPath:   envs.GetString("SECRET_MOUNT_PATH", "/etc/openfx"),
	}
}
