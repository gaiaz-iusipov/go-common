package httpservermw

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	httpheader "github.com/gaiaz-iusipov/go-common/http/header"
)

// BasicAuth implements a simple middleware for adding HTTP Basic Authentication.
func BasicAuth(username, password, realm string) func(next http.Handler) http.Handler {
	checkCredentials := func(reqUsername, reqPassword string) bool {
		return subtle.ConstantTimeCompare([]byte(username), []byte(reqUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(reqPassword)) == 1
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			reqUsername, reqPassword, ok := req.BasicAuth()
			if !ok || !checkCredentials(reqUsername, reqPassword) {
				rw.Header().Add(httpheader.WWWAuthenticate, fmt.Sprintf("Basic realm=%q", realm))
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(rw, req)
		})
	}
}
