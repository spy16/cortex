package server

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/chunked-app/cortex/core"
	"github.com/chunked-app/cortex/pkg/errors"
	"github.com/chunked-app/cortex/server/gql"
)

type Server struct {
	SystemInfo map[string]interface{}
	CoreAPI    *core.API
}

func (srv *Server) Serve(ctx context.Context, addr string) error {
	router := chi.NewRouter()
	router.Use(
		requestLogger(),
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		}).Handler,
	)

	router.NotFound(notFoundHandler())
	router.MethodNotAllowed(methodNotAllowedHandler())

	router.Get("/ping", pingHandler(srv.SystemInfo))
	router.Handle("/gql/play", playground.Handler("GraphQL playground", "/gql/query"))
	router.Handle("/gql/query", authenticate()(setupGQL(srv.CoreAPI)))

	return gracefulServe(ctx, 5*time.Second, addr, router)
}

func methodNotAllowedHandler() http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		respondJSON(wr, http.StatusMethodNotAllowed,
			errors.ErrInvalid.WithMsgf("%s not allowed for %s", req.Method, req.URL.Path))
		return
	}
}

func notFoundHandler() http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		respondJSON(wr, http.StatusNotFound,
			errors.ErrNotFound.WithMsgf("endpoint '%s %s' not found", req.Method, req.URL.Path))
		return
	}
}

func pingHandler(info map[string]interface{}) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		respondJSON(wr, http.StatusOK, info)
	}
}

func setupGQL(api *core.API) http.Handler {
	// setup GraphQL schema and handlers.
	schema := gql.NewExecutableSchema(gql.Config{
		Resolvers: &gql.Resolver{
			UsersAPI:  api,
			ChunksAPI: api,
		},
	})
	gqlServer := handler.NewDefaultServer(schema)
	gqlServer.AddTransport(&transport.POST{})
	gqlServer.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		ge := graphql.DefaultErrorPresenter(ctx, err)
		if ge.Unwrap() != nil {
			err = ge.Unwrap()
		} else {
			err = ge
		}

		e := errors.E(err)
		ext := map[string]interface{}{"code": e.Code}
		if e.Cause != "" {
			ext["cause"] = e.Cause
		}
		return &gqlerror.Error{
			Message:    e.Message,
			Extensions: ext,
		}
	})

	return gqlServer
}
