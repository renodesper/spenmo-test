package service

import (
	"github.com/spf13/viper"
)

type (
	// HealthService ...
	HealthService interface {
		Check() string
	}

	// HealthSvc ...
	HealthSvc struct{}
)

// Check ...
func (us *HealthSvc) Check() string {
	version := viper.GetString("app.version")
	if version == "" {
		return "1.0"
	}

	return version
}

// NewHealthService creates health service
func NewHealthService() HealthService {
	return &HealthSvc{}
}
