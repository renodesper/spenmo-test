package middleware

import (
	"github.com/go-kit/kit/endpoint"
)

type (
	// Middlewares contains middlewares that will be executed before
	// and after calling the endpoint. Each After middleware
	// should be deferred.
	Middlewares struct {
		Before []endpoint.Middleware
		After  []endpoint.Middleware
	}
)

// Chain will chain Before and After middlewares and traverse them
// in the order they're declared. That is, the first middleware
// of each field will be called last.
func Chain(middlewares Middlewares) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		next = ChainBefore(middlewares.Before...)(next)
		next = ChainAfter(middlewares.After...)(next)
		return next
	}
}

// ChainBefore will chain middlewares before calling the endpoint.
// For middlewares that were deferred, use ChainAfter.
func ChainBefore(middlewares ...endpoint.Middleware) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		for i := len(middlewares) - 1; i >= 0; i-- { // reverse
			next = middlewares[i](next)
		}
		return next
	}
}

// ChainAfter will chain middlewares after calling the endpoint.
// Use this for middlewares that were deferred.
func ChainAfter(middlewares ...endpoint.Middleware) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		for i := 0; i < len(middlewares); i++ {
			next = middlewares[i](next)
		}
		return next
	}
}
