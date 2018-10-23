package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Envs map[string]string

func NewEnvs() Envs {

	envs := Envs{}
	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		envs[parts[0]] = parts[1]
	}
	return envs
}

func (e Envs) GetInt(key string, defaultValue int) int {
	res := defaultValue
	if v, exists := e[key]; exists {
		intVal, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("defulat value is used because the environment value is not an int type: default value - '%d', environment variable - '%s'", defaultValue, e[key])
		} else {
			res = intVal
		}
	}
	return res
}

func (e Envs) GetDuration(key string, defaultValue time.Duration) time.Duration {
	res := defaultValue
	if v, exists := e[key]; exists {
		durationVal, err := time.ParseDuration(v)
		if err != nil {
			log.Printf("defulat value is used because the environment value is not an time.Duration type: default value - '%d', environment variable - '%s'", defaultValue, e[key])
		} else {
			res = durationVal
		}
	}
	return res
}

func (e Envs) GetString(key string, defaultValue string) string {
	if v, exists := e[key]; exists && len(v) > 0 {
		return v
	}
	return defaultValue
}

func (e Envs) GetBool(key string, defaultValue bool) bool {
	if v, exists := e[key]; exists {
		if v == "true" {
			return true
		} else {
			return false
		}
	}
	return defaultValue
}
