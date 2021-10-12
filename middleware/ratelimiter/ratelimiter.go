package ratelimiter

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/juju/ratelimit"
	"gitlab.com/renodesper/spenmo-test/util/errors"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

func CreateMiddleware(log logger.Logger, tb *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (result interface{}, e error) {
			if tb.TakeAvailable(1) == 0 {
				return nil, errors.LimitExceeded
			}

			return next(ctx, request)
		}
	}
}
