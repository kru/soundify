package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func CombineMiddleware(mds ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mds) - 1; i >= 0; i-- {
			m := mds[i]
			next = m(next)
		}
		return next
	}
}
