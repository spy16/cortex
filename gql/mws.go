package gql

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type middleware func(http.Handler) http.Handler

func verifyToken(authn Authenticator) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			h.ServeHTTP(wr, req)
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
