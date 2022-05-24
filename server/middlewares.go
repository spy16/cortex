package server

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/chunked-app/cortex/core/user"
)

type middleware func(http.Handler) http.Handler

func authenticate() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			// TODO: authenticate request using bearer token.
			ctx := user.Into(req.Context(), user.User{
				ID: "spy16",
			})

			next.ServeHTTP(wr, req.WithContext(ctx))
		})
	}
}

func requestLogger() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			startedAt := time.Now()
			fields := map[string]interface{}{
				"method":      req.Method,
				"path":        req.URL.Path,
				"remote_addr": req.RemoteAddr,
			}
			log.WithContext(req.Context()).WithFields(fields).Debugf("handling request")

			next.ServeHTTP(wr, req)

			fields["latency"] = time.Since(startedAt)
			log.WithContext(req.Context()).WithFields(fields).Infof("request completed")
		})
	}
}
