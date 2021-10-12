package metric

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

func CreateMiddleware(log logger.Logger, requestCount *kitprometheus.Counter, requestLatency *kitprometheus.Summary) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {

		return func(ctx context.Context, request interface{}) (result interface{}, e error) {
			defer func(begin time.Time) {
				lvs := []string{"method", "count", "error", "false"}
				requestCount.With(lvs...).Add(1)
				requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
			}(time.Now())

			return next(ctx, request)
		}
	}
}
