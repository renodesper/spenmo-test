package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"gitlab.com/renodesper/spenmo-test/service"
)

type HealthCheckResponse struct {
	Version string `json:"version"`
}

// MakeHealthCheckEndpoint ...
func MakeHealthCheckEndpoint(svc service.HealthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		version := svc.Check()
		return HealthCheckResponse{
			Version: version,
		}, nil
	}
}
