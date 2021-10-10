package recover

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	errs "gitlab.com/renodesper/spenmo-test/util/errors"
	"gitlab.com/renodesper/spenmo-test/util/logger"
)

// CreateMiddleware ...
func CreateMiddleware(log logger.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (result interface{}, e error) {
			defer func() {
				if r := recover(); r != nil {
					msg, ok := r.(string)
					if !ok {
						err, ok := r.(error)
						if ok {
							msg = err.Error()
						} else {
							msg = "Unexpected panic"
						}
					}
					e = errs.UnexpectedPanic.Wrap(errors.New(msg))
				}
			}()

			return next(ctx, request)
		}
	}
}
