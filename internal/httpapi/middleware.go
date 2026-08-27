package httpapi

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

func Auth(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		providedToken := auth[len(prefix):]

		if subtle.ConstantTimeCompare([]byte(token), []byte(providedToken)) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
