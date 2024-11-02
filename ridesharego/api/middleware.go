package api

import (
	"context"
	"database/sql"
	"main/auth"
	"main/core"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := core.SetupLogging()
		log = log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Logger

		ctx := context.WithValue(r.Context(), core.CtxLog, log)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withMiddleware(handler func(http.ResponseWriter, *http.Request, *logrus.Logger, *sql.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value(core.CtxLog).(*logrus.Logger)
		db := r.Context().Value(core.CtxAuth).(*sql.DB)
		handler(w, r, log, db)
	}
}

func authMiddleware(role string, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		authorizationParts := strings.Split(authorization, " ")
		if len(authorizationParts) != 2 {
			http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}

		if authorizationParts[0] != "Bearer" {
			http.Error(w, "invalid Authorization type", http.StatusUnauthorized)
			return
		}

		token := authorizationParts[1]
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		details, err := auth.GetUserDetailsAndValidate(token, role)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), core.CtxAuth, details)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
