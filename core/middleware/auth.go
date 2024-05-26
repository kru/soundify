package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kru/soundify/core/database"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user *database.User
		var err error
		user, err = database.QueryToken(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), database.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
