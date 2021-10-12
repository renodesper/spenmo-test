package http

import (
	"context"
	"net/http"
	"time"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator"
	"github.com/go-zoo/bone"
	jsoniter "github.com/json-iterator/go"
	"github.com/juju/ratelimit"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"gitlab.com/renodesper/spenmo-test/endpoint"
	m "gitlab.com/renodesper/spenmo-test/middleware"
	"gitlab.com/renodesper/spenmo-test/middleware/metric"
	"gitlab.com/renodesper/spenmo-test/middleware/ratelimiter"
	"gitlab.com/renodesper/spenmo-test/middleware/recover"
	ctxUtil "gitlab.com/renodesper/spenmo-test/util/ctx"
	e "gitlab.com/renodesper/spenmo-test/util/error"
	"gitlab.com/renodesper/spenmo-test/util/errors"
	"gitlab.com/renodesper/spenmo-test/util/logger"
	resp "gitlab.com/renodesper/spenmo-test/util/response"
)

var (
	json     = jsoniter.ConfigCompatibleWithStandardLibrary
	validate *validator.Validate
)

// NewHTTPHandler ...
func NewHTTPHandler(endpoints endpoint.Set, log logger.Logger) http.Handler {
	r := bone.New()

	// NOTE: Will be executed on the HTTP request object before the request is decoded
	serverRequestOpts := []httptransport.RequestFunc{
		ctxUtil.ExtractRequestID,
	}

	// NOTE: Will be executed on the HTTP response writer after the endpoint is invoked, but before anything written to the client
	serverResponseOpts := []httptransport.ServerResponseFunc{
		httptransport.SetContentType("application/json"),
	}

	serverOpts := []httptransport.ServerOption{
		httptransport.ServerBefore(serverRequestOpts...),
		httptransport.ServerAfter(serverResponseOpts...),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// NOTE: Middlewares
	publicMiddlewares := prepareMiddleware(log)

	// NOTE: Routes
	r.NotFound(http.HandlerFunc(notFound))

	GetHealthCheckEndpoint := m.Chain(publicMiddlewares)(endpoints.GetHealthCheckEndpoint)
	r.Get("/health", httptransport.NewServer(GetHealthCheckEndpoint, decodeNothing, encodeResponse, serverOpts...))

	GetAllUsersEndpoint := m.Chain(publicMiddlewares)(endpoints.GetAllUsersEndpoint)
	r.Get("/users", httptransport.NewServer(GetAllUsersEndpoint, decodeGetAllUsersRequest, encodeResponse, serverOpts...))

	GetUserByIDEndpoint := m.Chain(publicMiddlewares)(endpoints.GetUserByIDEndpoint)
	r.Get("/users/:id", httptransport.NewServer(GetUserByIDEndpoint, decodeGetUserRequest, encodeResponse, serverOpts...))

	CreateUserEndpoint := m.Chain(publicMiddlewares)(endpoints.CreateUserEndpoint)
	r.Post("/users", httptransport.NewServer(CreateUserEndpoint, decodeCreateUserRequest, encodeResponse, serverOpts...))

	UpdateUserEndpoint := m.Chain(publicMiddlewares)(endpoints.UpdateUserEndpoint)
	r.Put("/users/:id", httptransport.NewServer(UpdateUserEndpoint, decodeUpdateUserRequest, encodeResponse, serverOpts...))

	DeleteUserByIDEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteUserByIDEndpoint)
	r.Delete("/users/:id", httptransport.NewServer(DeleteUserByIDEndpoint, decodeDeleteUserRequest, encodeResponse, serverOpts...))

	GetAllTeamsEndpoint := m.Chain(publicMiddlewares)(endpoints.GetAllTeamsEndpoint)
	r.Get("/teams", httptransport.NewServer(GetAllTeamsEndpoint, decodeGetAllTeamsRequest, encodeResponse, serverOpts...))

	GetTeamEndpoint := m.Chain(publicMiddlewares)(endpoints.GetTeamEndpoint)
	r.Get("/teams/:id", httptransport.NewServer(GetTeamEndpoint, decodeGetTeamRequest, encodeResponse, serverOpts...))

	CreateTeamEndpoint := m.Chain(publicMiddlewares)(endpoints.CreateTeamEndpoint)
	r.Post("/teams", httptransport.NewServer(CreateTeamEndpoint, decodeCreateTeamRequest, encodeResponse, serverOpts...))

	UpdateTeamEndpoint := m.Chain(publicMiddlewares)(endpoints.UpdateTeamEndpoint)
	r.Put("/teams/:id", httptransport.NewServer(UpdateTeamEndpoint, decodeUpdateTeamRequest, encodeResponse, serverOpts...))

	DeleteTeamEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteTeamEndpoint)
	r.Delete("/teams/:id", httptransport.NewServer(DeleteTeamEndpoint, decodeDeleteTeamRequest, encodeResponse, serverOpts...))

	GetTeamMembersEndpoint := m.Chain(publicMiddlewares)(endpoints.GetTeamMembersEndpoint)
	r.Get("/team-members", httptransport.NewServer(GetTeamMembersEndpoint, decodeGetTeamMembersRequest, encodeResponse, serverOpts...))

	CreateTeamMemberEndpoint := m.Chain(publicMiddlewares)(endpoints.CreateTeamMemberEndpoint)
	r.Post("/team-members", httptransport.NewServer(CreateTeamMemberEndpoint, decodeCreateTeamMemberRequest, encodeResponse, serverOpts...))

	DeleteTeamMemberEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteTeamMemberEndpoint)
	r.Delete("/team-members", httptransport.NewServer(DeleteTeamMemberEndpoint, decodeDeleteTeamMemberRequest, encodeResponse, serverOpts...))

	DeleteTeamMembersByTeamIDEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteTeamMembersByTeamIDEndpoint)
	r.Delete("/team-members/teams/:teamId", httptransport.NewServer(DeleteTeamMembersByTeamIDEndpoint, decodeDeleteTeamMembersByTeamIDRequest, encodeResponse, serverOpts...))

	DeleteTeamMembersByUserIDEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteTeamMembersByUserIDEndpoint)
	r.Delete("/team-members/users/:userId", httptransport.NewServer(DeleteTeamMembersByUserIDEndpoint, decodeDeleteTeamMembersByUserIDRequest, encodeResponse, serverOpts...))

	GetAllWalletsEndpoint := m.Chain(publicMiddlewares)(endpoints.GetAllWalletsEndpoint)
	r.Get("/wallets", httptransport.NewServer(GetAllWalletsEndpoint, decodeGetAllWalletsRequest, encodeResponse, serverOpts...))

	GetWalletEndpoint := m.Chain(publicMiddlewares)(endpoints.GetWalletEndpoint)
	r.Get("/wallets/:id", httptransport.NewServer(GetWalletEndpoint, decodeGetWalletRequest, encodeResponse, serverOpts...))

	CreateWalletEndpoint := m.Chain(publicMiddlewares)(endpoints.CreateWalletEndpoint)
	r.Post("/wallets", httptransport.NewServer(CreateWalletEndpoint, decodeCreateWalletRequest, encodeResponse, serverOpts...))

	UpdateWalletEndpoint := m.Chain(publicMiddlewares)(endpoints.UpdateWalletEndpoint)
	r.Put("/wallets/:id", httptransport.NewServer(UpdateWalletEndpoint, decodeUpdateWalletRequest, encodeResponse, serverOpts...))

	DeleteWalletByIDEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteWalletByIDEndpoint)
	r.Delete("/wallets/:id", httptransport.NewServer(DeleteWalletByIDEndpoint, decodeDeleteWalletByIDRequest, encodeResponse, serverOpts...))

	DeleteWalletsByTeamIDEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteWalletsByTeamIDEndpoint)
	r.Delete("/wallets/teams/:teamId", httptransport.NewServer(DeleteWalletsByTeamIDEndpoint, decodeDeleteWalletsByTeamIDRequest, encodeResponse, serverOpts...))

	DeleteWalletsByUserIDEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteWalletsByUserIDEndpoint)
	r.Delete("/wallets/users/:userId", httptransport.NewServer(DeleteWalletsByUserIDEndpoint, decodeDeleteWallestByUserIDRequest, encodeResponse, serverOpts...))

	GetAllCardsEndpoint := m.Chain(publicMiddlewares)(endpoints.GetAllCardsEndpoint)
	r.Get("/cards", httptransport.NewServer(GetAllCardsEndpoint, decodeGetAllCardsRequest, encodeResponse, serverOpts...))

	GetCardEndpoint := m.Chain(publicMiddlewares)(endpoints.GetCardEndpoint)
	r.Get("/cards/:id", httptransport.NewServer(GetCardEndpoint, decodeGetCardRequest, encodeResponse, serverOpts...))

	CreateCardEndpoint := m.Chain(publicMiddlewares)(endpoints.CreateCardEndpoint)
	r.Post("/cards", httptransport.NewServer(CreateCardEndpoint, decodeCreateCardRequest, encodeResponse, serverOpts...))

	UpdateCardEndpoint := m.Chain(publicMiddlewares)(endpoints.UpdateCardEndpoint)
	r.Put("/cards/:id", httptransport.NewServer(UpdateCardEndpoint, decodeUpdateCardRequest, encodeResponse, serverOpts...))

	DeleteCardByIDEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteCardByIDEndpoint)
	r.Delete("/cards/:id", httptransport.NewServer(DeleteCardByIDEndpoint, decodeDeleteCardByIDRequest, encodeResponse, serverOpts...))

	DeleteCardsByWalletIDEndpoint := m.Chain(publicMiddlewares)(endpoints.DeleteCardsByWalletIDEndpoint)
	r.Delete("/cards/wallets/:walletId", httptransport.NewServer(DeleteCardsByWalletIDEndpoint, decodeDeleteCardsByWalletIDRequest, encodeResponse, serverOpts...))

	// NOTE: Prometheus metrics endpoint
	r.Get("/metrics", promhttp.Handler())

	return r
}

// decodeNothing returns (nil, nil) as placeholder for httptransport.DecodeRequestFunc
func decodeNothing(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	requestID := ctxUtil.GetRequestID(ctx)
	return json.NewEncoder(w).Encode(&resp.SuccessResponse{
		Data: response,
		Meta: resp.PopulateMeta(requestID),
	})
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	er := err.(e.Error)
	if er.Status == 0 {
		er.Status = 500
	}

	if viper.GetString("app.env") == "production" {
		er = er.WithoutStackTrace()
	}

	requestID := ctxUtil.GetRequestID(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	_ = json.NewEncoder(w).Encode(&resp.ErrorResponse{
		Errors: []e.Error{er},
		Meta:   resp.PopulateMeta(requestID),
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	_ = json.NewEncoder(w).Encode(&resp.ErrorResponse{
		Errors: []e.Error{errors.StatusNotFound.WithoutStackTrace()},
		Meta:   resp.PopulateMeta(r.Header.Get("X-Request-Id")),
	})
}

func prepareMiddleware(log logger.Logger) m.Middlewares {
	// NOTE: Rate Limiter
	ratelimitBucket := ratelimit.NewBucket(1*time.Second, 5)

	// NOTE: Metric
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "spenmo",
		Subsystem: "test",
		Name:      "request_count",
		Help:      "Number of requests received",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "spenmo",
		Subsystem: "test",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds",
	}, fieldKeys)

	publicMiddlewares := m.Middlewares{
		Before: []kitendpoint.Middleware{
			recover.CreateMiddleware(log),
			ratelimiter.CreateMiddleware(log, ratelimitBucket),
			metric.CreateMiddleware(log, requestCount, requestLatency),
		},
		After: []kitendpoint.Middleware{},
	}

	return publicMiddlewares
}
